package stone_plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/utils"
	"github.com/zanecloud/stone/stone_plugin/tools"
	"github.com/zanecloud/stone/stone_plugin/volume"
)

var ( // volumeNameRegex ensures the name assigned for the volume is valid.
	// This name is used to create the bind directory, so we need to avoid characters that
	// would make the path to escape the root directory.
	volumeNameRegex = utils.RestrictedVolumeNamePattern

	oldVfsDir = filepath.Join("vfs", "dir")

	validOpts = map[string]bool{
		"diskId":    true,
		"fsType":    true,
		"mediaType": true,
		"size":      true,
		"ioClass":   true,
		"exclusive": true,
	}
)

type validationError struct {
	error
}

type OptsConfig struct {
	diskId    string
	fsType    string // ext4, xfs...
	mediaType string // ssd, hdd
	size      int64  // Byte
	ioClass   int64  // ioClass: ssd(6-10), hdd(1-5)
	exclusive bool
}

func restoreVolumes(diskInfos []*tools.DiskInfo) (map[string]volume.Volume, error) {
	const VolumeRootPathName = "stone_volume"

	vols := map[string]volume.Volume{}
	for _, diskInfo := range diskInfos {
		rootPath := filepath.Join(diskInfo.MountPoint, VolumeRootPathName)
		if _, err := os.Stat(rootPath); os.IsNotExist(err) {
			continue
		}

		files, _ := ioutil.ReadDir(rootPath)
		for _, fi := range files {
			if fi.IsDir() {
				volumePath := filepath.Join(rootPath, fi.Name())
				v, err := volume.Restore(volumePath)
				if err != nil {
					return nil, err
				}
				vols[v.GetName()] = v
			}
		}
	}
	return vols, nil

}

func (s *stonePlugin) validateName(name string) error {
	if !volumeNameRegex.MatchString(name) {
		return validationError{fmt.Errorf("%q includes invalid characters for a local volume name, only %q are allowed", name, utils.RestrictedNameChars)}
	}
	return nil
}

func validateOpts(opts map[string]string) error {
	for optKey := range opts {
		if _, exist := validOpts[optKey]; !exist {
			return validationError{fmt.Errorf("optKey %s is invalid opt", optKey)}
		}
	}
	return nil
}

func (s *stonePlugin) setOpts(opts map[string]string) (*OptsConfig, error) {
	var err error
	if opts == nil || len(opts) == 0 {
		logrus.Warnf("opts is empty, use default size: 10G, ioClass: 0, fsType: random, mediaType: random")
	}

	if err := validateOpts(opts); err != nil {
		return nil, err
	}

	var size int64
	if sizeOpts, exists := opts["size"]; exists {
		size, err = strconv.ParseInt(sizeOpts, 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unknown format of size %s", opts["size"]))
		}
	} else {
		size = 10
	}

	var ioClass int64
	if ioClassOpts, exists := opts["ioClass"]; exists {
		ioClass, err = strconv.ParseInt(ioClassOpts, 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("unknown format of ioClass %s", opts["ioClass"]))
		}
	} else {
		ioClass = 0
	}

	v, exclusive := opts["exclusive"]
	if v != "" {
		exclusive, _ = strconv.ParseBool(v)
	}

	return &OptsConfig{
		diskId:    opts["diskId"],
		fsType:    opts["fsType"],
		mediaType: opts["mediaType"],
		size:      size,
		ioClass:   ioClass,
		exclusive: exclusive,
	}, nil
}

type scoredDisk struct {
	*tools.DiskInfo
	score int64
}

type scoredDiskList []*scoredDisk

func (g scoredDiskList) Len() int           { return len(g) }
func (g scoredDiskList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g scoredDiskList) Less(i, j int) bool { return g[i].score < g[j].score }
func (g scoredDiskList) Sort()              { sort.Sort(g) }

func (s *stonePlugin) chooseDisk(reqOpts *OptsConfig) (*tools.DiskInfo, error) {
	diskInfos, err := tools.Collect()
	if err != nil {
		return nil, err
	}

	if reqOpts.diskId != "" {
		for _, diskInfo := range diskInfos {
			if diskInfo.Id == reqOpts.diskId {
				return diskInfo, nil
			}
		}
		return nil, errors.New(fmt.Sprintf("Not found the disk with input diskId %s", reqOpts.diskId))
	}

	candidates := []*tools.DiskInfo{}

	sizeScores := scoredDiskList{}
	ioClassScores := scoredDiskList{}
	//TODO if opts.mediaType is "", allocate hdd first

OUTER:
	for _, diskInfo := range diskInfos {
		// if type is nil, means all types fit
		if reqOpts.fsType != "" && diskInfo.FsType != reqOpts.fsType {
			continue
		}
		if reqOpts.mediaType != "" && diskInfo.MediaType != reqOpts.mediaType {
			continue
		}

		var usedSize int64 = 0
		var existIoClass int64 = 0
		// collect all volumes on this disk
		for _, volume := range s.volumes {
			if volume.GetDiskId() == diskInfo.Id {
				if reqOpts.exclusive || volume.IsExclusive() {
					continue OUTER
				}
				usedSize += volume.GetSize()
				existIoClass += volume.GetIoClass()
			}
		}

		if diskInfo.Size-usedSize < reqOpts.size {
			// size is not enough
			continue
		}

		sizeScores = append(sizeScores, &scoredDisk{
			DiskInfo: diskInfo,
			score:    diskInfo.Size - usedSize,
		})

		if diskInfo.MediaType == tools.HDD {
			ioClassScores = append(ioClassScores, &scoredDisk{
				DiskInfo: diskInfo,
				score:    existIoClass * 2,
			})
		} else {
			ioClassScores = append(ioClassScores, &scoredDisk{
				DiskInfo: diskInfo,
				score:    existIoClass,
			})
		}
		candidates = append(candidates, diskInfo)
	}

	// sort
	sort.Sort(sizeScores)
	sort.Sort(ioClassScores)

	// choose
	var minScore = 100000
	var selectedOne *tools.DiskInfo
	for _, candidate := range candidates {
		var totalScore = 0

		for i, d := range sizeScores {
			if candidate.Id == d.Id {
				totalScore += i
			}
		}
		for i, d := range ioClassScores {
			if candidate.Id == d.Id {
				if d.score != 0 {
					totalScore += i
				}
			}
		}
		if totalScore < minScore {
			minScore = totalScore
			selectedOne = candidate
		}
	}
	if selectedOne == nil {
		return nil, errors.New("sorry, no disk suit.")
	}
	return selectedOne, nil
}

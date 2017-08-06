package local_volume

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/zanecloud/stone/stone_plugin/tools"
)

const (
	VolumeType         = "local"
	VolumeRootPathName = "stone_volume"
	VolumeDataPathName = "_data"
)

// localVolume implements the Volume interface from the volume package and
// represents the volumes created by Root.
type localVolume struct {
	m sync.Mutex
	// unique name of the volume
	Name string
	// volumePath
	VolumePath string
	// dataPath is the path on the host where the data lives
	DataPath string
	// driverName is the name of the driver that created the volume.
	DriverName string
	// volumeType is the type if the volume
	VolumeType string
	// diskId is the ID of diskInfo
	DiskId string
	// size is the size of the volume
	Size int64
	// ioClass
	IoClass int64
	// exclusive
	Exclusive bool
	// active refcounts the active mounts
	Active activeMount
}

type activeMount struct {
	Count   uint64
	Mounted bool
}

func New(driverName, volumeName string, size int64, ioClass int64, isExclusive bool, diskInfo *tools.DiskInfo) (*localVolume, error) {
	var err error
	v := &localVolume{
		Name:       volumeName,
		DriverName: driverName,
		VolumeType: VolumeType,
		DiskId:     diskInfo.Id,
		Size:       size,
		IoClass:    ioClass,
		Exclusive:  isExclusive,
	}

	// set volume path and data path
	v.VolumePath = filepath.Join(diskInfo.MountPoint, VolumeRootPathName, volumeName)
	v.DataPath = filepath.Join(diskInfo.MountPoint, VolumeRootPathName, volumeName, VolumeDataPathName)

	// allocate volume
	if err = v.allocateVolumeOnDisk(); err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(filepath.Dir(v.VolumePath))
		}
	}()

	// set disk quota
	sizeStr := strconv.FormatInt(size, 10)
	quotaId, err := tools.GetNextQuatoId()
	if err != nil {
		return nil, err
	}

	err = tools.SetDiskQuota(v.DataPath, sizeStr+"B", int(quotaId))
	if err != nil {
		return nil, err
	}

	// record volume to disk
	if err = v.toDisk(v.VolumePath); err != nil {
		return nil, err
	}

	return v, nil
}

func Remove(volumePath string) error {
	realPath, err := filepath.EvalSymlinks(volumePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		realPath = filepath.Dir(volumePath)
	}

	if err := tools.RemovePath(realPath); err != nil {
		return err
	}
	return tools.RemovePath(filepath.Dir(volumePath))
}

func Restore(volumePath string) (*localVolume, error) {
	v, err := fromDisk(volumePath)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Name returns the name of the given Volume.
func (v *localVolume) GetName() string {
	return v.Name
}

// DriverName returns the driver that created the given Volume.
func (v *localVolume) GetDriverName() string {
	return v.DriverName
}

// Path returns the data location.
func (v *localVolume) GetPath() string {
	return v.DataPath
}

// DataPath returns the constructed path of this volume.
func (v *localVolume) GetVolumePath() string {
	return v.VolumePath
}

// DataPath returns the constructed path of this volume.
func (v *localVolume) GetDataPath() string {
	return v.DataPath
}

// DiskId returns the ID of diskInfo
func (v *localVolume) GetDiskId() string {
	return v.DiskId
}

// Size returns the size of volume
func (v *localVolume) GetSize() int64 {
	return v.Size
}

func (v *localVolume) GetIoClass() int64 {
	return v.IoClass
}

func (v *localVolume) IsExclusive() bool {
	return v.Exclusive
}

// Mount implements the localVolume interface, returning the data location.
func (v *localVolume) Mount() (string, error) {
	v.m.Lock()
	defer v.m.Unlock()
	if !v.Active.Mounted {
		v.Active.Mounted = true
	}
	v.Active.Count++
	return v.DataPath, nil
}

// Umount is for satisfying the localVolume interface and does not do anything in this driver.
func (v *localVolume) Unmount() error {
	v.m.Lock()
	defer v.m.Unlock()
	v.Active.Count--
	if v.Active.Count == 0 {
		v.Active.Mounted = false
	}
	return nil
}

func (v *localVolume) Status() map[string]interface{} {
	return nil
}

func fromDisk(pth string) (*localVolume, error) {
	filename := filepath.Join(pth, "volume.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.New("error restore from disk without volume.json")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var obj localVolume
	err = json.Unmarshal(data, &obj)
	if err != nil {
		logrus.Debugf("parse %s failure: %s", filename, err)
		return nil, err
	}
	//if v.Name != obj.Name {
	//	return fmt.Errorf("name error: %s != %s", v.Name, obj.Name)
	//}
	dataPath := filepath.Join(pth, VolumeDataPathName)
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return nil, errors.New("error restore from disk without _data path")
	}

	if obj.Name == "" || obj.VolumePath == "" || obj.DataPath == "" || obj.DiskId == "" {
		return nil, errors.New("error restore from disk with wrong format data")
	}
	return &obj, nil
}

func (v *localVolume) toDisk(pth string) error {
	v.m.Lock()
	defer v.m.Unlock()

	data, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(pth, "volume.json"), data, 0600); err != nil {
		return err
	}
	return nil
}

package local_volume

import (
	"encoding/json"
	"fmt"
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
	name string
	// volumePath
	volumePath string
	// dataPath is the path on the host where the data lives
	dataPath string
	// driverName is the name of the driver that created the volume.
	driverName string
	// volumeType is the type if the volume
	volumeType string
	// diskId is the ID of diskInfo
	diskId string
	// size is the size of the volume
	size int64
	// ioClass
	ioClass int64
	// exclusive
	exclusive bool
	// active refcounts the active mounts
	active activeMount
}

type activeMount struct {
	count   uint64
	mounted bool
}

func New(driverName, volumeName string, size int64, ioClass int64, isExclusive bool, diskInfo *tools.DiskInfo) (*localVolume, error) {
	var err error
	v := &localVolume{
		name:       volumeName,
		driverName: driverName,
		volumeType: VolumeType,
		diskId:     diskInfo.Id,
		size:       size,
		ioClass:    ioClass,
		exclusive:  isExclusive,
	}

	// set volume path and data path
	v.volumePath = filepath.Join(diskInfo.MountPoint, VolumeRootPathName, volumeName)
	v.dataPath = filepath.Join(diskInfo.MountPoint, VolumeRootPathName, volumeName, VolumeDataPathName)

	// allocate volume
	if err = v.allocateVolumeOnDisk(); err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(filepath.Dir(v.volumePath))
		}
	}()

	// set disk quota
	sizeStr := strconv.FormatInt(size, 10)
	quotaId, err := tools.GetNextQuatoId()
	if err != nil {
		return nil, err
	}

	err = tools.SetDiskQuota(v.dataPath, sizeStr+"B", int(quotaId))
	if err != nil {
		return nil, err
	}

	// record volume to disk
	if err = v.ToDisk(v.volumePath); err != nil {
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

// Name returns the name of the given Volume.
func (v *localVolume) Name() string {
	return v.name
}

// DriverName returns the driver that created the given Volume.
func (v *localVolume) DriverName() string {
	return v.driverName
}

// Path returns the data location.
func (v *localVolume) Path() string {
	return v.dataPath
}

// DataPath returns the constructed path of this volume.
func (v *localVolume) VolumePath() string {
	return v.volumePath
}

// DataPath returns the constructed path of this volume.
func (v *localVolume) DataPath() string {
	return v.dataPath
}

// DiskId returns the ID of diskInfo
func (v *localVolume) DiskId() string {
	return v.diskId
}

// Size returns the size of volume
func (v *localVolume) Size() int64 {
	return v.size
}

func (v *localVolume) IoClass() int64 {
	return v.ioClass
}

func (v *localVolume) IsExclusive() bool {
	return v.exclusive
}

// Mount implements the localVolume interface, returning the data location.
func (v *localVolume) Mount() (string, error) {
	v.m.Lock()
	defer v.m.Unlock()
	if !v.active.mounted {
		v.active.mounted = true
	}
	v.active.count++
	return v.dataPath, nil
}

// Umount is for satisfying the localVolume interface and does not do anything in this driver.
func (v *localVolume) Unmount() error {
	v.m.Lock()
	defer v.m.Unlock()
	v.active.count--
	if v.active.count == 0 {
		v.active.mounted = false
	}
	return nil
}

func (v *localVolume) Status() map[string]interface{} {
	return nil
}

func (v *localVolume) FromDisk(pth string) error {
	filename := filepath.Join(filepath.Dir(pth), "volume.json")
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}

	v.m.Lock()
	defer v.m.Unlock()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var obj localVolume
	err = json.Unmarshal(data, &obj)
	if err != nil {
		logrus.Debugf("parse %s failure: %s", filename, err)
		return err
	}
	if v.name != obj.name {
		return fmt.Errorf("name error: %s != %s", v.name, obj.name)
	}
	v = &obj
	return nil
}

func (v *localVolume) ToDisk(pth string) error {
	v.m.Lock()
	defer v.m.Unlock()

	data, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(filepath.Dir(pth), "volume.json"), data, 0600); err != nil {
		return err
	}
	return nil
}

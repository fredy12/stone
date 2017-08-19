package local_volume

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/zanecloud/stone/stone_plugin/quota"
	"github.com/zanecloud/stone/stone_plugin/tools"
	"time"
)

const (
	VolumeType         = "local"
	VolumeRootPathName = "stone_volume"
	VolumeDataPathName = "_data"
)

// LocalVolume implements the Volume interface from the volume package and
// represents the volumes created by Root.
type LocalVolume struct {
	m sync.Mutex
	// unique name of the volume
	Name string
	// BackingFs
	BackingFs string
	// BasePath
	BasePath string
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
	// quotaControl
	quotaControl quota.QuotaControl
	// active refcounts the active mounts
	Active activeMount
}

type activeMount struct {
	Count   uint64
	Mounted bool
}

func New(driverName, volumeName string, size int64, ioClass int64, isExclusive bool, diskInfo *tools.DiskInfo) (*LocalVolume, error) {
	var err error
	v := &LocalVolume{
		Name:       volumeName,
		DriverName: driverName,
		VolumeType: VolumeType,
		DiskId:     diskInfo.Id,
		Size:       size,
		IoClass:    ioClass,
		Exclusive:  isExclusive,
		BackingFs:  diskInfo.FsType,
	}

	// set base path, volume path and data path
	v.BasePath = filepath.Join(diskInfo.MountPoint, VolumeRootPathName)
	v.VolumePath = filepath.Join(diskInfo.MountPoint, VolumeRootPathName, volumeName)
	v.DataPath = filepath.Join(diskInfo.MountPoint, VolumeRootPathName, volumeName, VolumeDataPathName)

	// allocate volume path
	if err = v.allocateVolumeOnDisk(v.VolumePath); err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(v.VolumePath)
		}
	}()

	// record volume to disk before quota set
	if err = v.toDisk(v.VolumePath); err != nil {
		return nil, err
	}

	// set disk quota
	qc, err := quota.NewQuota(v.BasePath, v.BackingFs)
	if err != nil {
		return nil, err
	}
	v.quotaControl = qc
	q := &quota.Quota{
		Size: uint64(size),
	}
	if err := qc.SetQuota(v.VolumePath, q); err != nil {
		return nil, err
	}

	// after quota set, allocate data path on disk
	if err = v.allocateVolumeOnDisk(v.DataPath); err != nil {
		return nil, err
	}

	return v, nil
}

func Remove(vol *LocalVolume, force bool) error {
	// remove disk quota
	if vol.quotaControl != nil {
		if err := vol.quotaControl.RemoveQuota(vol.VolumePath); err != nil {
			return err
		}
	}

	// remove path
	if err := tools.RemovePath(vol.VolumePath, force); err != nil {
		return err
	}

	return nil
}

func Restore(volumePath string) (*LocalVolume, error) {
	v, err := fromDisk(volumePath)
	if err != nil {
		return nil, err
	}
	// restore disk quota
	qc, err := quota.NewQuota(v.VolumePath, v.BackingFs)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("new quota control error when restore volume: %s", v.Name))
	}
	v.quotaControl = qc
	return v, nil
}

// Name returns the name of the given Volume.
func (v *LocalVolume) GetName() string {
	return v.Name
}

// DriverName returns the driver that created the given Volume.
func (v *LocalVolume) GetDriverName() string {
	return v.DriverName
}

// Path returns the data location.
func (v *LocalVolume) GetPath() string {
	return v.DataPath
}

// DataPath returns the constructed path of this volume.
func (v *LocalVolume) GetVolumePath() string {
	return v.VolumePath
}

// DataPath returns the constructed path of this volume.
func (v *LocalVolume) GetDataPath() string {
	return v.DataPath
}

// DiskId returns the ID of diskInfo
func (v *LocalVolume) GetDiskId() string {
	return v.DiskId
}

// Size returns the size of volume
func (v *LocalVolume) GetSize() int64 {
	return v.Size
}

func (v *LocalVolume) GetIoClass() int64 {
	return v.IoClass
}

func (v *LocalVolume) IsExclusive() bool {
	return v.Exclusive
}

// Mount implements the LocalVolume interface, returning the data location.
func (v *LocalVolume) Mount() (string, error) {
	v.m.Lock()
	defer v.m.Unlock()
	if !v.Active.Mounted {
		v.Active.Mounted = true
	}
	v.Active.Count++
	return v.DataPath, nil
}

// Umount is for satisfying the LocalVolume interface and does not do anything in this driver.
func (v *LocalVolume) Unmount() error {
	v.m.Lock()
	defer v.m.Unlock()
	v.Active.Count--
	if v.Active.Count == 0 {
		v.Active.Mounted = false
	}
	return nil
}

func (v *LocalVolume) Status() map[string]interface{} {
	time.Sleep(time.Second)
	return nil
}

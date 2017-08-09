package volume

import (
	"github.com/zanecloud/stone/stone_plugin/tools"
	"github.com/zanecloud/stone/stone_plugin/volume/local_volume"
)

type Volume interface {
	GetName() string
	GetPath() string
	GetVolumePath() string
	GetDiskId() string
	GetSize() int64
	GetIoClass() int64
	IsExclusive() bool
	Mount() (string, error)
	Unmount() error
	Status() map[string]interface{}
}

func New(driverName, volumeName string, size int64, ioClass int64, exclusive bool, diskInfo *tools.DiskInfo) (Volume, error) {
	return local_volume.New(driverName, volumeName, size, ioClass, exclusive, diskInfo)
}

func Remove(v Volume, force bool) error {
	return local_volume.Remove(v.(*local_volume.LocalVolume), force)
}

func Restore(volumePath string) (Volume, error) {
	return local_volume.Restore(volumePath)
}

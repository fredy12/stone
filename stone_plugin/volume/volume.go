package volume

import (
	"github.com/zanecloud/stone/stone_plugin/tools"
	"github.com/zanecloud/stone/stone_plugin/volume/local_volume"
)

type Volume interface {
	Name() string
	DriverName() string
	Path() string
	VolumePath() string
	DataPath() string
	DiskId() string
	Size() int64
	IoClass() int64
	IsExclusive() bool
	Mount() (string, error)
	Unmount() error
	Status() map[string]interface{}
	FromDisk(pth string) error
	ToDisk(pth string) error
}

func New(driverName, volumeName string, size int64, ioClass int64, exclusive bool, diskInfo *tools.DiskInfo) (Volume, error) {
	return local_volume.New(driverName, volumeName, size, ioClass, exclusive, diskInfo)
}

func Remove(v Volume) error {
	return local_volume.Remove(v.VolumePath())
}

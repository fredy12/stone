package system

// DiskInfo contains disk statistics of the host system
type DiskInfo struct {
	// filesystem name
	Filesystem string
	// The number of KBytes for a partition, 0 means unlimited
	TotalSize int64
	// The number of available KBytes for a partition, 0 means unlimited
	AvailSize int64
	// The number of used KBytes for a partition, 0 means unlimited
	UsedSize int64
	// Mount point on host
	MountedOn string
}

type DisksInfo []*DiskInfo

func (dsi DisksInfo) DiskInfo(mountedOn string) *DiskInfo {
	for _, di := range dsi {
		if di.MountedOn == mountedOn {
			return di
		}
	}
	return nil
}

package tools

const (
	SSD = "SSD"
	HDD = "HDD"
)

type DiskInfo struct {
	Id          string
	FsType      string
	MediaType   string
	FileSystem  string
	MountPoint  string
	Size        int64
	IoClass     int64
	IsBootDisk  bool
	IsExclusive bool
	Used        int64
	UsedPercent int64
}

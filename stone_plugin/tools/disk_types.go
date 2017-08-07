package tools

const (
	SSD = "SSD"
	HDD = "HDD"
)

type DiskInfo struct {
	Id         string
	FsType     string
	MediaType  string
	FileSystem string
	MountPoint string
	Size       int64
	IsBootDisk bool
}

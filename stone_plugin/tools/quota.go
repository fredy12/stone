package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
)

const (
	KILOBYTE   = 1024
	MEGABYTE   = 1024 * KILOBYTE
	GIGABYTE   = 1024 * MEGABYTE
	TERABYTE   = 1024 * GIGABYTE
	quotaMinId = uint32(16777216)
)

var (
	lock         sync.Mutex
	quotaLastId  uint32
	UseQuota                    = true
	quotaIds                    = make(map[uint32]uint32)
	mountPoints                 = make(map[uint64]string)
	bytesPattern *regexp.Regexp = regexp.MustCompile(`(?i)^(-?\d+)([KMGT]B?|B)$`)
)

func QuotaDriverStart(dir string) (string, error) {
	if !UseQuota {
		return "", nil
	}

	devId, err := GetDevId(dir)
	if err != nil {
		return "", err
	}

	lock.Lock()
	defer lock.Unlock()

	if mp, ok := mountPoints[devId]; ok {
		return mp, nil
	}

	mountPoint, hasQuota, _ := CheckMountpoint(devId)
	if len(mountPoint) == 0 {
		return mountPoint, fmt.Errorf("mountPoint not found: %s", dir)
	}
	if !hasQuota {
		doCmd("mount", "-o", "remount,grpquota", mountPoint)
	}

	vfsVersion, quotaFilename, err := GetVFSVersionAndQuotaFile(devId)
	if err != nil {
		return "", err
	}

	filename := mountPoint + "/" + quotaFilename
	if _, err := os.Stat(filename); err != nil {
		os.Remove(mountPoint + "/aquota.user")

		header := []byte{0x27, 0x19, 0xc0, 0xd9, 0x00, 0x00, 0x00, 0x00, 0x80, 0x3a, 0x09, 0x00, 0x80,
			0x3a, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x05, 0x00, 0x00, 0x00}
		if vfsVersion == "vfsv1" {
			header[4] = 0x01
		}

		if writeErr := ioutil.WriteFile(filename, header, 644); writeErr != nil {
			logrus.Errorf("write file error. %s, %s, %s", filename, vfsVersion, writeErr)
			return mountPoint, writeErr
		}

		if _, err := doCmd("setquota", "-g", "-t", "43200", "43200", mountPoint); err != nil {
			os.Remove(filename)
			return mountPoint, err
		}
		if err := setUserQuota(0, 0, mountPoint); err != nil {
			os.Remove(filename)
			return mountPoint, err
		}
	}

	// on
	out, err := doCmd("quotaon", "-pg", mountPoint)
	if strings.Contains(out, " is on") {
		mountPoints[devId] = mountPoint
		return mountPoint, nil
	}
	if _, err = doCmd("quotaon", mountPoint); err != nil {
		mountPoint = ""
	}

	mountPoints[devId] = mountPoint
	return mountPoint, err
}

//setfattr -n system.subtree -v $QUOTAID
func SetSubtree(dir string, qid uint32) (uint32, error) {
	if !UseQuota {
		return 0, nil
	}

	id := qid
	var err error
	if id == 0 {
		id = GetFileAttr(dir)
		if id > 0 {
			return id, nil
		}
		id, err = GetNextQuatoId()
	}

	if err != nil {
		return 0, err
	}
	strid := strconv.FormatUint(uint64(id), 10)
	_, err = doCmd("setfattr", "-n", "system.subtree", "-v", strid, dir)
	return id, err
}

func SetDiskQuota(dir string, size string, quotaId int) error {
	if !UseQuota {
		return nil
	}

	mountPoint, err := QuotaDriverStart(dir)
	if err != nil {
		return err
	}
	if len(mountPoint) == 0 {
		return fmt.Errorf("mountpoint not found: %s", dir)
	}

	id, err := SetSubtree(dir, uint32(quotaId))
	if id == 0 {
		return fmt.Errorf("subtree not found: %s %v", dir, err)
	}

	limit := toByteSize(size)
	return setUserQuota(id, limit, mountPoint)
}

func GetDevId(dir string) (uint64, error) {
	var st syscall.Stat_t
	if err := syscall.Stat(dir, &st); err != nil {
		logrus.Warnf("getDirDev: %s, %v", dir, err)
		return 0, err
	}
	return uint64(st.Dev), nil
}

func CheckMountpoint(devId uint64) (string, bool, string) {
	output, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		logrus.Warnf("ReadFile: %v", err)
		return "", false, ""
	}

	var mountPoint, fsType string
	hasQuota := false
	// /dev/sdf1 /apsarapangu/disk5 ext3 rw,noatime,nodiratime,errors=continue,barrier=1,data=ordered,grpquota 0 0
	for _, line := range strings.Split(string(output), "\n") {
		parts := strings.Split(line, " ")
		if len(parts) != 6 {
			continue
		}
		devId2, _ := GetDevId(parts[1])
		if devId == devId2 {
			mountPoint = parts[1]
			fsType = parts[2]
			for _, opt := range strings.Split(parts[3], ",") {
				if opt == "grpquota" {
					hasQuota = true
				}
			}
			break
		}
	}

	return mountPoint, hasQuota, fsType
}

func doCmd(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		logrus.Warnf("cmd: [%s %s], output: %s err: %v", name, strings.Join(args, " "), string(output), err)
	} else {
		logrus.Infof("cmd: [%s %s]", name, strings.Join(args, " "))
	}
	return string(output), err
}

func setUserQuota(quotaId uint32, diskQuota uint64, mountPoint string) error {
	uid := strconv.FormatUint(uint64(quotaId), 10)
	limit := strconv.FormatUint(diskQuota, 10)
	_, err := doCmd("setquota", "-g", uid, "0", limit, "0", "0", mountPoint)
	return err
}

//getfattr -n system.subtree --only-values --absolute-names /
func GetFileAttr(dir string) uint32 {
	v := 0
	out, err := doCmd("getfattr", "-n", "system.subtree", "--only-values", "--absolute-names", dir)
	if err == nil {
		v, _ = strconv.Atoi(out)
	}
	return uint32(v)
}

func SetFileAttr(dir string, id uint32) error {
	strid := strconv.FormatUint(uint64(id), 10)
	_, err := doCmd("setfattr", "-n", "system.subtree", "-v", strid, dir)
	return err
}

//load
//repquota -gan
//Group           used    soft    hard  grace    used  soft  hard  grace
//----------------------------------------------------------------------
//#0        --  494472       0       0            938     0     0
//#54       --       8       0       0              2     0     0
//#4        --      16       0       0              4     0     0
//#22       --      28       0       0              4     0     0
//#16777220 +- 2048576       0 2048575              9     0     0
//#500      --   47504       0       0            101     0     0
//#16777221 -- 3048576       0 3048576              8     0     0
func loadQuotaIds() (uint32, error) {
	minId := quotaMinId
	output, err := doCmd("repquota", "-gan")
	if err != nil {
		return minId, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) == 0 || line[0] != '#' {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) == 0 {
			continue
		}
		id, err := strconv.Atoi(parts[0][1:])
		uid := uint32(id)
		if err == nil && uid > quotaMinId {
			quotaIds[uid] = 1
			if uid > minId {
				minId = uid
			}
		}
	}
	logrus.Infof("Load repquota ids: %d, list: %v", len(quotaIds), quotaIds)
	return minId, nil
}

//next id
func GetNextQuatoId() (uint32, error) {
	lock.Lock()
	defer lock.Unlock()

	if quotaLastId == 0 {
		var err error
		quotaLastId, err = loadQuotaIds()
		if err != nil {
			return 0, err
		}
	}
	id := quotaLastId
	for {
		if id < quotaMinId {
			id = quotaMinId
		}
		id++
		if _, ok := quotaIds[id]; !ok {
			break
		}
	}
	quotaIds[id] = 1
	quotaLastId = id
	return id, nil
}

func toByteSize(s string) uint64 {
	parts := bytesPattern.FindStringSubmatch(strings.TrimSpace(s))
	if len(parts) < 3 {
		return 0
	}

	value, err := strconv.ParseUint(parts[1], 10, 0)
	if err != nil || value < 1 {
		return 0
	}

	var bytes uint64
	unit := strings.ToUpper(parts[2])
	switch unit[:1] {
	case "T":
		bytes = value * TERABYTE
	case "G":
		bytes = value * GIGABYTE
	case "M":
		bytes = value * MEGABYTE
	case "K":
		bytes = value * KILOBYTE
	case "B":
		bytes = value
	}

	if bytes < KILOBYTE {
		return 1
	}

	return bytes / KILOBYTE
}

func GetVFSVersionAndQuotaFile(devId uint64) (string, string, error) {
	output, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		logrus.Warnf("ReadFile: %v", err)
		return "", "", err
	}

	vfsVersion := "vfsv0"
	quotaFilename := "aquota.group"
	// /dev/sdf1 /apsarapangu/disk5 ext3 rw,noatime,nodiratime,errors=continue,barrier=1,data=ordered,grpquota 0 0
	for _, line := range strings.Split(string(output), "\n") {
		parts := strings.Split(line, " ")
		if len(parts) != 6 {
			continue
		}

		devId2, _ := GetDevId(parts[1])
		if devId == devId2 {
			for _, opt := range strings.Split(parts[3], ",") {
				items := strings.SplitN(opt, "=", 2)
				if len(items) != 2 {
					continue
				}
				switch items[0] {
				case "jqfmt":
					vfsVersion = items[1]
				case "grpjquota":
					quotaFilename = items[1]
				}
			}
			break
		}
	}

	return vfsVersion, quotaFilename, nil
}

package main

/*

Filesystem      1K-blocks      Used Available Use% Mounted on
/dev/sda2        51475068   5446652  43390592  12% /
devtmpfs        131918316         0 131918316   0% /dev

*/

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"crypto/md5"
	"encoding/hex"
	log "github.com/Sirupsen/logrus"
)

const (
	SSD = "SSD"
	HDD = "HDD"
)

var (
	ErrDfNotResult    = errors.New("Command df not output result")
	ErrDfNotOkResult  = errors.New("Command df not right result")
	ErrIostatNoResult = errors.New("Command iostat not output result")
)

type DfInfo struct {
	Filesystem  string
	Type        string
	Block       int64
	Used        int64
	Available   int64
	UsedPercent int64
	Mounted     string
}

func ParseDf() ([]*DfInfo, error) {
	out, err := exec.Command("/bin/sh", "-c", "df -T -B 1 -P 2>/dev/null").Output()
	if err != nil {
		if len(out) == 0 {
			log.Errorf("df error %v", err)
			return nil, err
		}
	}

	lines := bytes.Split(out, []byte("\n"))

	if len(lines) == 0 || len(lines) == 1 {
		return nil, ErrDfNotResult
	}

	i := 0
	infos := make([]*DfInfo, len(lines))

	// skip the first title line.
	for _, l := range lines[1:] {
		fields := strings.Fields(string(l))
		if len(fields) != 7 {
			continue
		}

		df := &DfInfo{}

		df.Filesystem = fields[0]

		df.Type = fields[1]

		if block, err := strconv.ParseInt(fields[2], 10, 64); err != nil {
			log.Errorf("%v", err)
			continue
		} else {
			df.Block = block
		}

		if used, err := strconv.ParseInt(fields[3], 10, 64); err != nil {
			log.Errorf("%v", err)
			continue
		} else {
			df.Used = used
		}

		if avail, err := strconv.ParseInt(fields[4], 10, 64); err != nil {
			log.Errorf("%v", err)
			continue
		} else {
			df.Available = avail
		}

		if usedPercent, err := strconv.ParseInt(strings.Trim(fields[5], "%"), 10, 64); err != nil {
			log.Errorf("%v", err)
			continue
		} else {
			df.UsedPercent = usedPercent
		}

		df.Mounted = fields[6]

		infos[i] = df
		i++
	}

	if i == 0 {
		return nil, ErrDfNotOkResult
	}
	return infos[:i], nil
}

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

func Collect() ([]*DiskInfo, error) {
	deviceFlag := make(map[string]string)
	deviceFlagType := map[string]string{
		"0": SSD,
		"1": HDD,
	}

	diskDevices, err := Command("ls", "/sys/block")
	if err != nil {
		return nil, err
	}

	for _, device := range strings.Fields(diskDevices) {
		rotationPath := fmt.Sprintf("/sys/block/%s/queue/rotational", device)
		diskType, err := Command("cat", rotationPath)
		if err != nil {
			return nil, err
		}
		diskType = strings.TrimSpace(diskType)
		deviceFlag[device] = deviceFlagType[diskType]
	}

	cwd := fmt.Sprintf("/proc/%d/cwd", os.Getpid())
	bootPath, err := Command("readlink", "-f", cwd)
	if err != nil {
		return nil, err
	}

	// 通过 df 命令查询磁盘的空间信息
	infos, err := ParseDf()
	if err != nil {
		log.Errorf("get all disk message failed: %v", err)
		return nil, err
	}
	if len(infos) == 0 {
		return nil, err
	}

	diskInfos := make([]*DiskInfo, 0, len(infos))
	for _, info := range infos {
		if info.Mounted == bootPath {
			continue
		}

		if skipMountedPath(info.Mounted) {
			continue
		}

		if skipFilesystem(info.Filesystem) {
			continue
		}

		// skip all Filesystem not in /sys/block/
		diskType := "unknown"
		for k, v := range deviceFlag {
			if strings.Contains(info.Filesystem, k) {
				diskType = v
			}
		}
		if diskType == "unknown" {
			log.Warnf("Unknown disk type found for device: %s", info.Filesystem)
			continue
		}

		// NOTICE: check nvme
		if _, err := os.Stat(info.Filesystem); err == nil && strings.HasPrefix(info.Filesystem, "/dev/nvme") {
			shell := fmt.Sprintf("ll %s | awk -F, '{print $1}' | awk '{print $5}'", info.Filesystem)
			// 259 is device's primary num
			if RunShellInt(shell) == 259 {
				diskType = "nvme"
			}
		}

		diskInfos = append(diskInfos, &DiskInfo{
			Id:          Md5sum(fmt.Sprintf("%s_%s_%s", info.Filesystem, info.Type, diskType)),
			FileSystem:  info.Filesystem,
			Size:        info.Block,
			MountPoint:  info.Mounted,
			IsBootDisk:  false,
			FsType:      info.Type,
			MediaType:   diskType,
			Used:        info.Used,
			UsedPercent: info.UsedPercent,
		})
	}
	return diskInfos, nil
}

func Md5sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func RunShellInt(shell string) int {
	out, err := exec.Command("/bin/sh", "-c", shell).Output()
	if err != nil {
		log.Error("RunShellInt %s failed: %v", shell, err)
		return -1
	}

	n, err := strconv.Atoi(string(bytes.Trim(out, " \n\r\t ")))
	if err != nil {
		log.Error("RunShellInt %s failed: %v", shell, err)
		return -1
	}

	return n
}

func Command(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		log.Warnf("cmd: [%s %s], output: %s err: %v", name, strings.Join(args, " "), string(output), err)
	} else {
		log.Infof("cmd: [%s %s]", name, strings.Join(args, " "))
	}
	return string(output), err
}

func skipMountedPath(mount string) bool {
	mount = strings.TrimSpace(mount)

	// skip boot path
	if mount == "/" {
		return true
	}
	return false
}

func skipFilesystem(fs string) bool {
	fs = strings.TrimSpace(fs)

	// skip /dev/rbd0, /dev/rbd1 and so on.
	if !strings.HasPrefix(fs, "/dev/rbd") {
		return false
	}

	suffix := strings.TrimPrefix(fs, "/dev/rbd")
	if _, err := strconv.Atoi(suffix); err != nil {
		return false
	}

	return true
}

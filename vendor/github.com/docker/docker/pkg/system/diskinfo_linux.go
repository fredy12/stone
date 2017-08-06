package system

import (
	"os/exec"
	"strconv"
	"strings"
)

func ReadDisksInfo() (DisksInfo, error) {
	out, err := exec.Command("df").Output()
	if err != nil {
		return nil, err
	}

	var dsi DisksInfo
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.Fields(line)
		if len(parts) != 6 {
			continue
		}

		di := &DiskInfo{}
		di.Filesystem = parts[0]
		di.MountedOn = parts[5]

		di.TotalSize, err = strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			return nil, err
		}
		di.UsedSize, err = strconv.ParseInt(parts[2], 10, 0)
		if err != nil {
			return nil, err
		}
		di.AvailSize, err = strconv.ParseInt(parts[3], 10, 0)
		if err != nil {
			return nil, err
		}

		dsi = append(dsi, di)
	}

	return dsi, nil
}

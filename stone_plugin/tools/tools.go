package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

func Collect() ([]*DiskInfo, error) {
	var DiskCmd = "docker-disk"
	args := []string{
		"list",
	}
	out, err := exec.Command(DiskCmd, args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %q", err.Error(), out)
	}
	var diskInfoS []*DiskInfo
	if err := json.Unmarshal(out, &diskInfoS); err != nil {
		return nil, err
	}
	return diskInfoS, nil
}

func Command(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		logrus.Warnf("cmd: [%s %s], output: %s err: %v", name, strings.Join(args, " "), string(output), err)
	} else {
		logrus.Infof("cmd: [%s %s]", name, strings.Join(args, " "))
	}
	return string(output), err
}

func RemovePath(path string) error {
	if err := os.RemoveAll(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

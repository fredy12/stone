package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

func Collect(useRootDisk bool) ([]*DiskInfo, error) {
	var DiskCmd = "docker-disk"
	args := []string{
		"list",
	}
	if useRootDisk {
		args = append(args, "rootdisk")
	}
	out, err := exec.Command(DiskCmd, args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %q", err.Error(), out)
	}
	var diskInfos []*DiskInfo
	if err := json.Unmarshal(out, &diskInfos); err != nil {
		return nil, err
	}
	return diskInfos, nil
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

func RemovePath(path string, force bool) error {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		if err := os.RemoveAll(path); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
	} else {
		if !force {
			return errors.New(fmt.Sprintf("%s: no such file or directory", path))
		}
	}

	return nil
}

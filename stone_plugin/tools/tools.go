package tools

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"crypto/md5"
	"github.com/Sirupsen/logrus"
)

func Command(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		logrus.Warnf("cmd: [%s %s], output: %s err: %v", name, strings.Join(args, " "), string(output), err)
	} else {
		logrus.Infof("cmd: [%s %s]", name, strings.Join(args, " "))
	}
	return string(output), err
}

func RunShellInt(shell string) int {
	out, err := exec.Command("/bin/sh", "-c", shell).Output()
	if err != nil {
		logrus.Error("RunShellInt %s failed: %v", shell, err)
		return -1
	}

	n, err := strconv.Atoi(string(bytes.Trim(out, " \n\r\t ")))
	if err != nil {
		logrus.Error("RunShellInt %s failed: %v", shell, err)
		return -1
	}

	return n
}

func Md5sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}

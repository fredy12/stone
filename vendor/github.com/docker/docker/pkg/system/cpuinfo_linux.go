package system

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadCpuInfo() (*CpuInfo, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseCpuInfo(file)
}

func parseCpuInfo(reader io.Reader) (*CpuInfo, error) {
	var (
		cpuinfo          CpuInfo
		err              error
		currentSocket    int64
		currentCore      int64
		currentProcessor int64
	)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}

		switch parts[0] {
		case "physical id":
			currentSocket, err = strconv.ParseInt(parts[1], 10, 16)
			if err != nil {
				return nil, err
			}
		case "core id":
			currentCore, err = strconv.ParseInt(parts[1], 10, 16)
			if err != nil {
				return nil, err
			}
		case "processor":
			currentProcessor, err = strconv.ParseInt(parts[1], 10, 16)
			if err != nil {
				return nil, err
			}
			cpuinfo.Processors += 1
		case "":
			cpuinfo.AddSocket(currentSocket)
			socket := cpuinfo.Sockets.Get(currentSocket)
			socket.AddCore(currentCore)
			socket.Cores.Get(currentCore).AddProcessor(currentProcessor)
		}
	}
	// Handle errors that may have occurred during the reading of the file.
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &cpuinfo, nil
}

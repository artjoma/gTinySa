package main

import (
	"errors"
	"fmt"
	"go.bug.st/serial/enumerator"
	"strconv"
	"strings"
)

// chmod 666 /dev/ttyACM0
func tryFindDevice() (*enumerator.PortDetails, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, errors.New("no serial ports found")
	}

	for _, port := range ports {
		vid, err := strconv.ParseInt(port.VID, 16, 64)
		if err != nil {
			continue
		}
		pid, err := strconv.ParseInt(port.PID, 16, 64)
		if err != nil {
			continue
		}
		if pid == tinyDevicePid && tinyDeviceVid == vid {
			fmt.Printf("TinySa detected pid:%d vid:%d name:%s \n", pid, vid, port.Name)
			return port, nil
		}
	}

	return nil, nil
}

// TinySa return command inside body response
func removeCommandLine(response string) []string {
	data := strings.Split(strings.TrimSpace(response), "\n")
	if len(data) == 0 {
		return []string{}
	}
	return data[1:]
}

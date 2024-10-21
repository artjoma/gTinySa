package main

import (
	"bytes"
	"fmt"
	"go.bug.st/serial"
	"strconv"
	"strings"
)

const (
	tinyDeviceVid = int64(0x0483)
	tinyDevicePid = int64(0x5740)
	scanPoints    = uint16(400)
)

var (
	DefaultTinyConfig = &TinyConfig{
		scanPoints: scanPoints,
	}
)

// USB API https://tinysa.org/wiki/pmwiki.php?n=Main.USBInterface
type TinySaDevice struct {
	port   serial.Port
	config *TinyConfig
}

type TinyConfig struct {
	scanPoints uint16
}

func NewTinyDevice(deviceName string, config *TinyConfig) (*TinySaDevice, error) {
	port, err := serial.Open(deviceName, &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
	})
	if err != nil {
		return nil, err
	}

	return &TinySaDevice{port: port, config: config}, nil
}

func (tiny *TinySaDevice) Version() ([]string, error) {
	data, err := tiny.version()
	if err != nil {
		return nil, err
	}
	return removeCommandLine(data), nil
}

func (tiny *TinySaDevice) version() (string, error) {
	if err := tiny.sendRequest("version"); err != nil {
		return "", err
	}
	return tiny.readResponse()
}

func (tiny *TinySaDevice) Scan(fromFq, toFq uint64) ([][]int64, error) {
	dataStr, err := tiny.scan(fromFq, toFq, tiny.config.scanPoints)
	if err != nil {
		return nil, err
	}
	dataArr := removeCommandLine(dataStr)
	scanResult := make([][]int64, 0, len(dataArr))
	for _, record := range dataArr {
		cols := strings.Split(record, " ")
		freq, err := strconv.ParseInt(cols[0], 10, 64)
		if err != nil {
			return nil, err
		}
		level, err := strconv.ParseFloat(cols[1], 10)
		levelI := int64(0)
		if err == nil {
			levelI = int64(level)
		} else {
			fmt.Println(record)
			levelI = scanResult[len(scanResult)-1][1]
		}
		scanResult = append(scanResult, []int64{freq, levelI})
	}

	return scanResult, nil
}

// scan 90M 100M 100 3
func (tiny *TinySaDevice) scan(fromFq, toFq uint64, points uint16) (string, error) {
	if err := tiny.sendRequest(fmt.Sprintf("scan %d %d %d 3", fromFq, toFq, points)); err != nil {
		return "", err
	}
	return tiny.readResponse()
}

func (tiny *TinySaDevice) sendRequest(request string) error {
	_, err := tiny.port.Write([]byte(request + "\n\r"))
	return err
}

func (tiny *TinySaDevice) Close() {
	tiny.port.Close()
}

func (tiny *TinySaDevice) readResponse() (string, error) {
	charBuff := make([]byte, 1)
	line := ""
	buf := bytes.NewBufferString("")
	for {
		_, err := tiny.port.Read(charBuff)
		if err != nil {
			return "", err
			break
		}
		c := charBuff[0]
		if c == 13 { // \r
			continue
		}
		line += string(c)
		if c == 10 { // \n
			buf.WriteString(line)
			line = ""
			continue
		}
		// command end
		if strings.HasSuffix(line, "ch>") {
			break
		}
	}

	return buf.String(), nil
}

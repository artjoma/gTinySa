package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	portDetails, err := tryFindDevice()
	if err != nil {
		fmt.Println("Err:" + err.Error())
		return
	}
	if portDetails == nil {
		fmt.Println("No device found")
		return
	}

	device, err := NewTinyDevice(portDetails.Name, DefaultTinyConfig)
	if err != nil {
		log.Fatal(err)
	}

	// self check,
	for i := 0; i < 8; i++ {
		device.Version()
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Println(device.Version())
	fmt.Println(device.Info())

	values, err := device.Scan(90_000_000, 105_000_000)
	fmt.Println(values, err)
	device.port.Close()
}

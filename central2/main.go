package main

import (
	"fmt"
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	serviceUUID = bluetooth.NewUUID([16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3})
	charUUID    = bluetooth.NewUUID([16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3})
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return err
	}

	// Scan
	var foundDevice bluetooth.ScanResult
	fmt.Printf("Scanning...\n")
	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.LocalName() != "tinygo ble peripheral" {
			return
		}
		foundDevice = result

		// Stop the scan.
		err := adapter.StopScan()
		if err != nil {
			// Unlikely, but we can't recover from this.
			println("failed to stop the scan:", err.Error())
		}
	})
	if err != nil {
		return err
	}

	// Found a peripheral. Connect to it.
	fmt.Printf("Connecting to %q (%s)\n", foundDevice.LocalName(), foundDevice.Address.String())
	device, err := adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}

	// Connected.
	fmt.Printf("Discovering service...\n")
	services, err := device.DiscoverServices([]bluetooth.UUID{serviceUUID})
	if err != nil {
		return err
	}
	service := services[0]

	// Get characteristics present in this service.
	fmt.Printf("Get characteristics...\n")
	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{charUUID})
	if err != nil {
		return err
	}

	fmt.Printf("running\n")
	tick := time.Tick(500 * time.Millisecond)
	buf := []byte{0x00}
	for {
		<-tick
		buf[0]++
		chars[0].WriteWithoutResponse(buf)
	}

	return nil
}

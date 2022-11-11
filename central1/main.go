package main

import (
	"fmt"
	"log"

	"tinygo.org/x/bluetooth"
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

	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		fmt.Printf("%#v %s\n", result.Address.String(), result.LocalName())
		return
	})

	return nil
}

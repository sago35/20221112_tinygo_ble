package main

import (
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

	adv := adapter.DefaultAdvertisement()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "tinygo ble peripheral",
	})
	if err != nil {
		return err
	}

	err = adv.Start()
	if err != nil {
		return err
	}

	return nil
}

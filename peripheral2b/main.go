package main

import (
	"image/color"
	"log"
	"machine"
	"time"

	"tinygo.org/x/bluetooth"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/gophers"
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

func disp() error {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
		SCL:       machine.SCL0_PIN,
		SDA:       machine.SDA0_PIN,
	})
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: ssd1306.Address_128_32,
		Width:   128,
		Height:  32,
	})
	display.ClearDisplay()

	font := &gophers.Regular32pt
	//str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str := "ABEFGIJMNVWXYZ" + "ABEFG"
	tick := time.Tick(200 * time.Millisecond)
	i := 0
	for {
		<-tick

		display.ClearBuffer()
		tinyfont.WriteLine(&display, font, 5, 30, str[i:i+5], color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
		i = (i + 1) % (len(str) - 5)
		display.Display()

	}

	return nil
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

	ch := make(chan byte, 10)
	err = adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID: charUUID,
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission |
					bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) != 1 {
						return
					}
					ch <- value[0]
				},
			},
		},
	})
	if err != nil {
		return err
	}

	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
		SCL:       machine.SCL0_PIN,
		SDA:       machine.SDA0_PIN,
	})
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: ssd1306.Address_128_32,
		Width:   128,
		Height:  32,
	})
	display.ClearDisplay()

	machine.LED_GREEN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	i := 0
	font := &gophers.Regular32pt
	str := "ABEFGIJMNVWXYZ" + "ABEFG"
	for {
		val := <-ch
		if (val % 2) == 0 {
			machine.LED_GREEN.Low()
		} else {
			machine.LED_GREEN.High()
		}

		display.ClearBuffer()
		tinyfont.WriteLine(&display, font, 5, 30, str[i:i+5], color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
		i = (i + 1) % (len(str) - 5)
		display.Display()
	}

	return nil
}

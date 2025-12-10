package main

import (
	_ "embed"
	"image/color"
	"machine"
	"time"

	"m4-apps/lib/ntpclient"
	"m4-apps/lib/rgb75"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

func main() {

	for !machine.Serial.DTR() {
		time.Sleep(100 * time.Millisecond)
	}

	ntpclient.SyncSystemTime("0.pool.ntp.org:123")

	device := rgb75.NewDevice()

	colors := []color.RGBA{
		{255, 0, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
		{0, 255, 255, 255},
		{0, 0, 255, 255},
		{255, 0, 255, 255},
		{255, 255, 255, 255},
	}

	for {

		pst := time.FixedZone("PST", -8*60*60)
		t := time.Now().In(pst)

		tinyfont.WriteLine(device, &proggy.TinySZ8pt7b, 15, 12, t.Format("Jan 2"), colors[4])
		tinyfont.WriteLineColors(device, &proggy.TinySZ8pt7b, 2, 26, t.Format("3:04:05pm"), colors)

		device.Display()

		time.Sleep(1 * time.Millisecond)
	}

}

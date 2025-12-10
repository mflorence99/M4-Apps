package main

import (
	"fmt"
	"image/color"
	"m4-apps/lib/rgb75"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

func main() {

	device := rgb75.NewDevice()

	tick := time.Tick(1 * time.Second)

	for t := range tick {

		tinyfont.WriteLine(device, &freesans.Regular9pt7b, 2, 20, fmt.Sprintf("%d", t.Second()), color.RGBA{255, 128, 255, 255})

		device.Display()

		time.Sleep(1 * time.Millisecond)
	}

}

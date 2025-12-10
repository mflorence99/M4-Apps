package main

import (
	"image/color"
	"m4-apps/lib/rgb75"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

func main() {

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
		tinyfont.WriteLine(device, &freesans.Regular9pt7b, 12, 14, "Hello", colors[4])
		tinyfont.WriteLineColors(device, &freesans.Regular9pt7b, 2, 28, "Buster!", colors)

		device.Display()

		time.Sleep(10 * time.Microsecond)
	}

}

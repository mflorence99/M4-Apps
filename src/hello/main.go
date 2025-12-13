package main

import (
	"m4-apps/lib/rgb75"
	"m4-apps/lib/utils"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

func main() {

	device := rgb75.NewDevice()

	for {
		tinyfont.WriteLine(device, &freesans.Regular9pt7b, 12, 14, "Hello", utils.Colors[4])
		tinyfont.WriteLineColors(device, &freesans.Regular9pt7b, 2, 28, "Buster!", utils.Colors)

		device.Display()

		time.Sleep(10 * time.Microsecond)
	}

}

package main

import (
	"image/color"
	"m4-apps/rgb75"
	"machine"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

func main() {

	h75 := rgb75.New(
		machine.HUB75_OE, machine.HUB75_LAT, machine.HUB75_CLK,
		[6]machine.Pin{
			machine.HUB75_R1, machine.HUB75_G1, machine.HUB75_B1,
			machine.HUB75_R2, machine.HUB75_G2, machine.HUB75_B2,
		},
		[]machine.Pin{
			machine.HUB75_ADDR_A, machine.HUB75_ADDR_B, machine.HUB75_ADDR_C,
			machine.HUB75_ADDR_D, machine.HUB75_ADDR_E,
		})

	err := h75.Configure(rgb75.Config{
		Width:      64,
		Height:     32,
		ColorDepth: 4,
		DoubleBuf:  true,
	})

	if err != nil {
		for {
			println("error: " + err.Error())
			time.Sleep(time.Second)
		}
	}

	h75.Resume()

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

		tinyfont.WriteLine(h75, &freesans.Regular9pt7b, 12, 14, "Hello", colors[4])
		tinyfont.WriteLineColors(h75, &freesans.Regular9pt7b, 2, 28, "Buster!", colors)

		h75.Display()
		time.Sleep(10 * time.Microsecond)
	}

}

package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

func main() {
	// Configure the NeoPixel pin (D4) as an output pin.
	// The TinyGo machine package provides a direct alias for the pin.
	neoPixelPin := machine.PA23
	neoPixelPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Create a new WS2812 driver instance.
	// The Matrix Portal M4 has a single NeoPixel.
	led := ws2812.New(neoPixelPin)

	// Define the colors to use. The NeoPixel on this board uses GRB (Green, Red, Blue) ordering.
	// We use the standard color.RGBA struct and let the driver handle the color order.
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 0} // Off state

	for {
		// Set the NeoPixel color to orange and show it.
		// The Write command sends the color data to the LED.
		led.WriteColors([]color.RGBA{red})
		time.Sleep(1 * time.Second)

		// Turn the NeoPixel off.
		led.WriteColors([]color.RGBA{black})
		time.Sleep(1 * time.Second)
	}
}

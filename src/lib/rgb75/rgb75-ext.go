package rgb75

import (
	"machine"
)

// 🟧 My extensions to the original driver

func NewDevice() *Device {

	device := New(
		machine.HUB75_OE, machine.HUB75_LAT, machine.HUB75_CLK,
		[6]machine.Pin{
			machine.HUB75_R1, machine.HUB75_G1, machine.HUB75_B1,
			machine.HUB75_R2, machine.HUB75_G2, machine.HUB75_B2,
		},
		[]machine.Pin{
			machine.HUB75_ADDR_A, machine.HUB75_ADDR_B, machine.HUB75_ADDR_C,
			machine.HUB75_ADDR_D, machine.HUB75_ADDR_E,
		})

	device.Configure(Config{
		Width:      64,
		Height:     32,
		ColorDepth: 4,
		DoubleBuf:  true,
	})

	device.Resume()

	return device

}

//go:build debug

package utils

import (
	"machine"
	"time"
)

// 🟧 Wait for the serial port to be opened for debugging

func WaitForSerial() {
	for !machine.Serial.DTR() {
		time.Sleep(100 * time.Millisecond)
	}
}

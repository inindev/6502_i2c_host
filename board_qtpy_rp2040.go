//go:build qtpy_rp2040

package main

import "machine"

const (
	NEOPIXEL     = machine.GPIO12
	NEOPIXEL_PWR = machine.GPIO11
	LED          = machine.NoPin

	INT_PIN = machine.GPIO3
	INT_LED = LED

	I2C_SDA       = machine.I2C1_QT_SDA_PIN
	I2C_SCL       = machine.I2C1_QT_SCL_PIN
	I2C_Frequency = 1e6
)

var (
	I2C_PORT = machine.I2C1
)

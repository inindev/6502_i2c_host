//go:build feather_rp2040

package main

import "machine"

const (
	NEOPIXEL     = machine.GPIO16
	NEOPIXEL_PWR = machine.NoPin
	LED          = machine.GPIO13

	INT_PIN = machine.GPIO6

	I2C_SDA       = machine.I2C1_SDA_PIN
	I2C_SCL       = machine.I2C1_SCL_PIN
	I2C_Frequency = 1e6
)

var (
	I2C_PORT = machine.I2C1
)

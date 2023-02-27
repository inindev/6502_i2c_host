package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

type Neopixel struct {
	ws ws2812.Device
}

func (np Neopixel) flashError(code uint8) {
	for i := uint8(0); i < code; i++ {
		switch {
		case ((i + 1) % 8) == 0:
			np.flashGreen(500)
		case ((i + 1) % 4) == 0:
			np.flashBlue(500)
		default:
			np.flashRed(250)
		}
	}
}

func (np Neopixel) flashRed(delayMs time.Duration) {
	np.flash(0x200000, delayMs)
}

func (np Neopixel) flashGreen(delayMs time.Duration) {
	np.flash(0x001000, delayMs)
}

func (np Neopixel) flashBlue(delayMs time.Duration) {
	np.flash(0x000020, delayMs)
}

func (np Neopixel) flash(color uint32, delayMs time.Duration) {
	np.ws.WriteByte(uint8(color >> 8))  // g
	np.ws.WriteByte(uint8(color >> 16)) // r
	np.ws.WriteByte(uint8(color))       // b
	time.Sleep(time.Millisecond * delayMs)

	np.ws.WriteByte(0) // g
	np.ws.WriteByte(0) // r
	np.ws.WriteByte(0) // b
	time.Sleep(time.Millisecond * delayMs)
}

func NewNeopixel(neo, neoPwr machine.Pin) Neopixel {
	if neoPwr != machine.NoPin {
		neoPwr.Configure(machine.PinConfig{Mode: machine.PinOutput})
		neoPwr.High()
	}

	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ws := ws2812.New(neo)
	ws.WriteByte(0x00)
	time.Sleep(time.Microsecond * 300)

	return Neopixel{ws}
}

package main

import (
	"machine"
	"time"
)

var code = []uint8{
	0xa9, 0x42, //        lda #$42
	0x8d, 0x10, 0x80, //  sta $8010
}

var (
	ram             [0x10000]uint8 // 64k ram
	i2cGpio         I2cGpio
	inputEventTicks int64
)

func main() {
	i2cGpio = configIO()

	// store program at address
	prog(0x6000)

	for {
		base := 0
		for i := base; i < base+0x100; i++ {
			x := ram[i]
			y := revBits(x)

			// write unit_20, banks a & b
			val := uint16(x)<<8 | uint16(y)
			err := i2cGpio.setValue(UNIT_20, val)
			if err != nil {
				doPanic(err, 5)
			}

			// write unit_21, bank b only
			err = i2cGpio.setValueB(UNIT_21, y)
			if err != nil {
				doPanic(err, 6)
			}

			if inputEventTicks > 0 {
				ticks := inputEventTicks
				inputEventTicks = 0 // reset event

				val, err := i2cGpio.getValueA(UNIT_21) // read bank a, unit_21
				if err != nil {
					doPanic(err, 7)
				}

				processInputEvent(val, ticks)
			}

			time.Sleep(time.Millisecond * 50)

			if base == 0 && i == 255 {
				time.Sleep(time.Second)
			}
		}
	}
}

// store code at addr
func prog(addr uint16) {
	copy(ram[0:], code)

	// nop all memory
	for i := 0; i < len(ram); i++ {
		ram[i] = 0xea
	}

	// demo: counter 0-255
	for i := 0; i < 256; i++ {
		ram[i] = uint8(i)
	}

	// demo: sweep 256-511
	for i := 256; i < 512; i++ {
		m16 := i % 16
		if m16 < 8 {
			ram[i] = uint8(1 << m16)
		} else {
			ram[i] = uint8(1 << (15 - m16))
		}
	}

	// copy code to addr
	for i, v := range code {
		ram[addr+uint16(i)] = v
	}

	// store execution address in reset vector
	ram[0xffec] = uint8(addr)
	ram[0xffed] = uint8(addr << 8)
}

func onInterrupt(pin machine.Pin) {
	inputEventTicks = time.Now().UnixMilli()
}

func configIO() I2cGpio {
	// configure interrupt pin
	INT_PIN.Configure(machine.PinConfig{Mode: machine.PinInput | machine.PinInputPulldown})
	INT_PIN.SetInterrupt(machine.PinFalling, onInterrupt)

	i2cGpio, err := NewI2cGpio(I2C_PORT, I2C_SDA, I2C_SCL)
	if err != nil {
		doPanic(err, 1)
	}

	err = i2cGpio.setupIO(UNIT_20, 0x0000) // unit_20: all output
	if err != nil {
		doPanic(err, 2)
	}

	err = i2cGpio.setupIO(UNIT_21, 0xff00) // unit_21: bank a input, b output
	if err != nil {
		doPanic(err, 3)
	}

	err = i2cGpio.setupInterrupt(UNIT_21, 0xff00) // unit_21: bank a interrupt on change
	if err != nil {
		doPanic(err, 4)
	}

	return i2cGpio
}

func doPanic(err error, num uint8) {
	var neopixel Neopixel
	switch {
	case NEOPIXEL != machine.NoPin:
		neopixel = NewNeopixel(NEOPIXEL, NEOPIXEL_PWR)
	case LED != machine.NoPin:
		LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	for {
		print("*** panic ***  (error ")
		print(num)
		print("): ")
		println(err)

		switch {
		case NEOPIXEL != machine.NoPin:
			neopixel.flashError(num)
		case LED != machine.NoPin:
			for i := uint8(0); i < num; i++ {
				switch {
				case ((i + 1) % 4) == 0:
					LED.Set(true)
					time.Sleep(time.Millisecond * 500)
					LED.Set(false)
					time.Sleep(time.Millisecond * 500)
				default:
					LED.Set(true)
					time.Sleep(time.Millisecond * 250)
					LED.Set(false)
					time.Sleep(time.Millisecond * 250)
				}
			}
		}

		time.Sleep(time.Second * 3)
	}
}

func toHex(x uint8) string {
	const D16 = "0123456789abcdef"
	return string([]uint8{D16[x>>4], D16[x&0xf]})
}

func revBits(x uint8) uint8 {
	x = (x&0xaa)>>1 | (x&0x55)<<1
	x = (x&0xcc)>>2 | (x&0x33)<<2
	x = (x&0xf0)>>4 | (x&0x0f)<<4
	return x
}

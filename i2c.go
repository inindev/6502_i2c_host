package main

import (
	"machine"
)

const (
	UNIT_20 = uint8(0x20)
	UNIT_21 = uint8(0x21)

	REG_IODIR   = uint8(0x00) // i/o direction
	REG_IPOL    = uint8(0x02) // input polarity
	REG_GPINTEN = uint8(0x04) // interrupt-on-change
	REG_GPPUA   = uint8(0x0c) // pull-up resistor
	REG_GPIO    = uint8(0x12) // output bufue
)

type I2cGpio struct {
	port *machine.I2C
}

var buf = []uint8{0, 0}

func (ig I2cGpio) getValue(id uint8) (uint16, error) {
	buf[0] = 0
	buf[1] = 0
	err := ig.port.ReadRegister(id, REG_GPIO, buf)
	if err != nil {
		return 0, err
	}

	val := uint16(buf[0])<<8 | uint16(buf[1])
	return val, nil
}

func (ig I2cGpio) setValue(id uint8, val uint16) error {
	buf[0] = uint8(val >> 8) // bank a value
	buf[1] = uint8(val)      // bank b value
	return ig.port.WriteRegister(id, REG_GPIO, buf)
}

func (ig I2cGpio) getValueA(id uint8) (uint8, error) {
	return ig.getValueN(id, 0x00)
}

func (ig I2cGpio) setValueA(id uint8, val uint8) error {
	return ig.setValueN(id, 0x00, val)
}

func (ig I2cGpio) getValueB(id uint8) (uint8, error) {
	return ig.getValueN(id, 0x01)
}

func (ig I2cGpio) setValueB(id uint8, val uint8) error {
	return ig.setValueN(id, 0x01, val)
}

func (ig I2cGpio) getValueN(id, bank uint8) (uint8, error) {
	buf[0] = 0
	buf[1] = 0
	err := ig.port.ReadRegister(id, REG_GPIO|bank, buf[:1])
	if err != nil {
		return 0, err
	}

	var val = buf[0]
	return val, nil
}

func (ig I2cGpio) setValueN(id, bank uint8, val uint8) error {
	buf[0] = val // bank value
	buf[1] = 0
	return ig.port.WriteRegister(id, REG_GPIO|bank, buf[:1])
}

func (ig I2cGpio) setupIO(id uint8, iomask uint16) error {
	buf[0] = uint8(iomask >> 8) // bank a pullup
	buf[1] = uint8(iomask)      // bank b pullup
	err := ig.port.WriteRegister(id, REG_GPPUA, buf)
	if err != nil {
		return err
	}

	buf[0] = uint8(iomask >> 8) // bank a out/in direction
	buf[1] = uint8(iomask)      // bank b out/in direction
	return ig.port.WriteRegister(id, REG_IODIR, buf)
}

func (ig I2cGpio) setupInterrupt(id uint8, intmask uint16) error {
	// The GPINTEN register controls the interrupt-on-change feature for each pin.
	// If a bit is set, the corresponding pin is enabled for interrupt-on-change.
	// The DEFVAL and INTCON registers must also be configured if any pins are enabled for interrupt-on-change.
	buf[0] = uint8(intmask >> 8) // bank a interrupt enable
	buf[1] = uint8(intmask)      // bank b interrupt enable
	return ig.port.WriteRegister(UNIT_21, REG_GPINTEN, buf)
}

func NewI2cGpio(port *machine.I2C, sda, scl machine.Pin) (I2cGpio, error) {
	err := port.Configure(machine.I2CConfig{
		Frequency: 1000 * machine.KHz,
		SDA:       sda,
		SCL:       scl,
	})

	if err != nil {
		return I2cGpio{}, err
	}

	return I2cGpio{machine.I2C1}, nil
}

package main

type Address struct {
	i2c       I2cGpio
	i2cBankId uint8
}

type Data struct {
	i2c       I2cGpio
	i2cBankId uint8
	bank      uint8
}

type Gpio struct {
	i2c       I2cGpio
	i2cBankId uint8
	bank      uint8
}

func (a Address) getValue() (uint16, error) {
	return a.i2c.getValue(a.i2cBankId)
}

func NewAddress(i2c I2cGpio, i2cBankId uint8) Address {
	return Address{
		i2c:       i2c,
		i2cBankId: i2cBankId,
	}
}

func (d Data) getValue() (uint8, error) {
	return d.i2c.getValueN(d.i2cBankId, d.bank)
}

func NewData(i2c I2cGpio, i2cBankId, bank uint8) Data {
	return Data{
		i2c:       i2c,
		i2cBankId: i2cBankId,
		bank:      bank,
	}
}

const (
	oRST     = 0
	oCLK     = 1
	oLED     = 2
	iBTN_RST = 3
	iBTN_CLK = 4
	iRWB     = 5
)

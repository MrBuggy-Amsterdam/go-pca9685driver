package pca9685

import (
	"encoding/binary"
	"errors"

	log "github.com/MrBuggy-Amsterdam/buggy-log"
	smbus "github.com/corrupt/go-smbus"
)

type Register struct {
	size            int
	i2caddress      uint8 // For example, the pca has hardware address 0x40
	registerAddress uint8 // For example, the mode 1 register has register address 0x0
	bus             *smbus.SMBus
}

func (reg *Register) Read() (int, error) {
	buf := make([]byte, reg.size)

	_, err := reg.bus.Read_block_data(reg.registerAddress, buf)
	if err != nil {
		return -99, err
	}

	if reg.size == 1 {
		return int(buf[0]), nil
	} else if reg.size == 2 {
		resData := binary.LittleEndian.Uint16(buf)
		return int(resData), nil
	} else {
		return -99, errors.New("invalid register size")
	}
}

func (reg *Register) Write(data []byte) error {
	if len(data) != reg.size {
		log.Error("Tried to write %d bytes to %d-sized register", len(data), reg.size)
		return errors.New("data size mismatched with register size")
	}

	if len(data) <= 0 {
		return errors.New("tried to write 0 bytes")
	}

	var err error
	if reg.size == 1 {
		err = reg.bus.Write_byte_data(reg.registerAddress, data[0])
	} else if reg.size == 2 || reg.size == 4 {
		for i := 0; i < reg.size; i++ {
			err = reg.bus.Write_byte_data(reg.registerAddress+uint8(i), data[i])
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("invalid register size")
	}
	return err
}

func newRegister(bus *smbus.SMBus, size int, i2caddress uint8, registerAddress uint8) *Register {
	return &Register{size, i2caddress, registerAddress, bus}
}

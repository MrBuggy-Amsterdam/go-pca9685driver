package pca9685

import (
	"encoding/binary"
	"errors"
)

type PwmChannel struct {
	register *Register
}

func (ch *PwmChannel) SetDutyCycle(value uint16) error {
	if value > 0xFFFF {
		return errors.New("invalid value as duty cycle")
	}

	data := make([]byte, 4)
	if value == 0xFFFF {
		binary.LittleEndian.PutUint16(data[2:], 0x1000)
		return ch.register.Write(data)
	} else {
		// Shift our value by four because the PCA9685 is only 12 bits but our value is 16
		actualValue := (value + 1) >> 4

		binary.LittleEndian.PutUint16(data[2:], actualValue)
		return ch.register.Write(data)
	}
}

func newPwmChannel(register *Register) *PwmChannel {
	return &PwmChannel{register}
}

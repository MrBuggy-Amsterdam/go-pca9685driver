package pca9685

import (
	"time"

	log "github.com/MrBuggy-Amsterdam/buggy-log"
	smbus "github.com/corrupt/go-smbus"
)

type PCA9685 struct {
	bus *smbus.SMBus

	// Special registers
	mode1register    *Register
	mode2register    *Register
	prescaleRegister *Register

	// PWM registers (tied to a channel)
	pwmRegisters []Register
	PwmChannels  []PwmChannel
}

const (
	referenceClockSpeed float32 = 25000000
)

// Initializes the PCA board, use i2cdetect to find the address on the bus you are using
// by default, the address is 0x40
// the suggested frequency is 50Hz (for servos)
func Initialize(busNr uint8, address uint8, frequency uint8) *PCA9685 {
	log.Debug("Initializing PCA9685 controller")

	bus, err := smbus.New(uint(busNr), byte(address))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Prepare special registers
	mode1reg := newRegister(bus, 1, address, 0x00)
	mode2reg := newRegister(bus, 1, address, 0x01)
	prescaleRegister := newRegister(bus, 1, address, 0xFE)

	// Prepare array of PWM registers and channels
	pwmRegisters := make([]Register, 16)
	pwmChannels := make([]PwmChannel, 16)

	// Prepare PWM registers and channels and add them to the array
	for i := uint8(0); i < 16; i++ {
		newReg := newRegister(bus, 4, address, 0x06+(i*4))
		pwmRegisters[i] = *newReg
		pwmChannels[i] = *newPwmChannel(newReg)
	}

	pca := &PCA9685{bus, mode1reg, mode2reg, prescaleRegister, pwmRegisters, pwmChannels}

	log.Debug("Initialized PCA9685 controller on bus %d (addr: 0x%x)", bus, address)
	pca.reset()
	pca.setFrequency(float32(frequency))

	return pca
}

func (pca *PCA9685) reset() {
	err := pca.bus.Write_byte(0x0) // mode 1
	if err != nil {
		log.Fatal("Reset error: %s", err.Error())
	}
}

// Set the frequency in hertz
func (pca *PCA9685) setFrequency(frequency float32) {
	prescale := uint8(referenceClockSpeed/4096.0/frequency + 0.5)
	log.Debug("Calculated prescale %d", prescale)

	if prescale < 3 {
		log.Fatal("Cannot configure PCA9685 at the given frequency")
	}

	oldMode := 0x00

	// reusable buffer for writing
	writeBytes := make([]byte, 1)

	// mode 1, sleep
	writeBytes[0] = byte((oldMode & 0x7F) | 0x10)
	pca.mode1register.Write(writeBytes)

	// prescale
	writeBytes[0] = byte(prescale)
	pca.prescaleRegister.Write(writeBytes)

	// mode 1
	writeBytes[0] = byte(oldMode)
	pca.mode1register.Write(writeBytes)

	time.Sleep(1 * time.Second)

	// mode 1, autoincrement on, fix to stop pca9685 from accepting commands at all addresses
	writeBytes[0] = byte(oldMode | 0x40)
	pca.mode1register.Write(writeBytes)

	log.Debug("prescale: %v", prescale)
}

# go-pca9685driver

This is a Go module useful for controlling the pca9685 16-channel PWM board on a Raspberry Pi or similar devices. The board is controlled over i2c. This implementation is based on [the Adafruit Python3 implementation](https://github.com/adafruit/Adafruit_CircuitPython_PCA9685).

> **note**: this driver is by no means done or fully tested. We are working on a more extensive driver with more flexible options.

## Installation

```bash
go get github.com/MrBuggy-Amsterdam/go-pca9685driver
```

## Usage

After initializing the board on a specific **bus** (the i2c device used), **address** (0x40 by default) and **frequency** (50Hz is preferred for servo control). You will have access to a `PwmChannel` which corresponds to one of the 16 channels on the board. On these channels you can set a duty cycle between 0 (0% of the period is a high voltage) and 65535 (100% of the period is a high voltage). 

```Go
package main

import (
    pca9685 "github.com/MrBuggy-Amsterdam/go-pca9685driver"
)

func main() {
	// Use i2cdetect to find the address on the bus you are using
	pca := pca9685.Initialize(5, 0x40, 50)

	// Get the first PWM channel (left-most one) to control
	firstChannel := pca.PwmChannels[0]
	// A value of 4900 means that the signal will be high for 4900/65535 = 0.075% of the time
	// with a frequency of 50 Hz (period is 20ms), this means that the signal will be high for about 0.075 * 20 = 1.5ms
	firstChannel.SetDutyCycle(4900)

	secondChannel := pca.PwmChannels[1]
	secondChannel.SetDutyCycle(0)
}
```

## Troubleshooting

A good start is to use `i2cdetect` to scan your i2c bus and confirm that your device is connected. Also make sure to read the [datasheet](https://www.nxp.com/products/power-management/lighting-driver-and-controller-ics/led-controllers/16-channel-12-bit-pwm-fm-plus-ic-bus-led-controller:PCA9685) and familiarize yourself with PWM control.

Open an issue or PR to let us know of any issues you encounter.
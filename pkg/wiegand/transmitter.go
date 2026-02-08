package wiegand

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	// PinD0         = 17
	// PinD1         = 27
	// PulseDuration = 80 * time.Microsecond
	// BitInterval   = 4 * time.Millisecond
)

// Transmit отправляет бинарную последовательность на контроллер Wiegand
func Transmit(binaryCode string) error {
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()

	pin0 := rpio.Pin(PinD0)
	pin1 := rpio.Pin(PinD1)

	pin0.Input()
	pin1.Input()

	for _, bit := range binaryCode {
		var pin rpio.Pin
		if bit == '0' {
			pin = pin0
		} else {
			pin = pin1
		}

		pin.Output()
		pin.Low()
		time.Sleep(PulseDuration)
		pin.Input()

		time.Sleep(BitInterval)
	}

	return nil
}
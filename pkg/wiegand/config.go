package wiegand

import "time"

const (
	PinD0         = 17
	PinD1         = 27
	Timeout       = 200 * time.Millisecond
	MinBits       = 24
	PulseDuration = 80 * time.Microsecond
	BitInterval   = 4 * time.Millisecond
)

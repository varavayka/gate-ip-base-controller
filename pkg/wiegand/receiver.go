package wiegand

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)


var (
	bits     string
	lastTick time.Time
	mu       sync.Mutex
)
var binaryPresentCard string

func Receiver() string {

	// Открываем доступ к GPIO
	if err := rpio.Open(); err != nil {
		log.Fatalf("Ошибка открытия GPIO: %v", err)
	}
	defer rpio.Close()

	pin0 := rpio.Pin(PinD0)
	pin1 := rpio.Pin(PinD1)

	// Настраиваем пины: вход с подтяжкой вверх и детекцией падающего фронта
	pin0.Input()
	pin0.PullUp()
	pin0.Detect(rpio.FallEdge)

	pin1.Input()
	pin1.PullUp()
	pin1.Detect(rpio.FallEdge)

	// Обработчики в горутинах
	go func() {
		for {
			if pin0.EdgeDetected() {
				mu.Lock()
				bits += "0"
				lastTick = time.Now()
				mu.Unlock()
			}
			time.Sleep(50 * time.Microsecond)
		}
	}()

	go func() {
		for {
			if pin1.EdgeDetected() {
				mu.Lock()
				bits += "1"
				lastTick = time.Now()
				mu.Unlock()
			}
			time.Sleep(50 * time.Microsecond)
		}
	}()

	// Основной цикл проверки таймаута
	for {
		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		if bits != "" && time.Since(lastTick) > Timeout {
			if len(bits) >= MinBits {
				binaryPresentCard += processCard(bits)
			}
			bits = ""
			break
		}
		mu.Unlock()
	}
	return binaryPresentCard
}

func processCard(binaryData string) string {

	var fullValue uint64
	for _, bit := range binaryData {
		fullValue <<= 1
		if bit == '1' {
			fullValue |= 1
		}
	}


	if len(binaryData) == 26 {
		payload := binaryData[1:25]

		facilityBits := payload[0:8]
		var facilityCode uint64
		for _, bit := range facilityBits {
			facilityCode <<= 1
			if bit == '1' {
				facilityCode |= 1
			}
		}

		cardBits := payload[8:24]
		var cardNumber uint64
		for _, bit := range cardBits {
			cardNumber <<= 1
			if bit == '1' {
				cardNumber |= 1
			}
		}

		

	}  else {
		fmt.Printf("Формат: Нестандартный (%d бит)\n", len(binaryData))
	}
	return binaryData
}

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

func Receiver(resultChan chan<- string) {
	if err := rpio.Open(); err != nil {
		log.Printf("Ошибка открытия GPIO: %v", err)
		return
	}
	// Закроется только в самом конце работы Receiver
	defer rpio.Close()

	pin0 := rpio.Pin(PinD0)
	pin1 := rpio.Pin(PinD1)

	pin0.Input()
	pin0.PullUp()
	pin0.Detect(rpio.FallEdge)

	pin1.Input()
	pin1.PullUp()
	pin1.Detect(rpio.FallEdge)

	// Канал для остановки фоновых горутин
	done := make(chan struct{})

	// Очищаем старые данные перед началом
	mu.Lock()
	bits = ""
	mu.Unlock()

	// Воркер для опроса пинов
	worker := func(p rpio.Pin, bit string) {
		for {
			select {
			case <-done: // Если получили сигнал об окончании — выходим из горутины
				return
			default:
				if p.EdgeDetected() {
					mu.Lock()
					bits += bit
					lastTick = time.Now()
					mu.Unlock()
				}
				time.Sleep(50 * time.Microsecond)
			}
		}
	}

	go worker(pin0, "0")
	go worker(pin1, "1")

	var finalCard string

	// Основной цикл проверки таймаута
	for {
		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		if bits != "" && time.Since(lastTick) > Timeout {
			if len(bits) >= MinBits {
				finalCard = processCard(bits)
			}
			bits = ""
			mu.Unlock()
			break // Выходим из цикла ожидания
		}
		mu.Unlock()
	}

	// 1. Сначала останавливаем горутины
	close(done)
	
	// 2. Даем им микропаузу, чтобы они успели завершиться до rpio.Close()
	time.Sleep(10 * time.Millisecond)

	// 3. Отправляем результат
	resultChan <- finalCard
	// close(resultChan) // Закрывай канал здесь, если уверен, что больше никто в него не пишет
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

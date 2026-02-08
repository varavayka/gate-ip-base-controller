package wiegand

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/stianeikeland/go-rpio/v4"
)



var (
	pin0 rpio.Pin
	pin1 rpio.Pin
)

func Jammer() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("Ошибка открытия GPIO: %v", err)
	}
	defer cleanup()

	pin0 = rpio.Pin(PinD0)
	pin1 = rpio.Pin(PinD1)

	// Обработка Ctrl+C для корректного завершения
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("\n\nПрерывание пользователем...")
		cleanup()
		os.Exit(0)
	}()

	// Изначально разблокируем
	jamOff()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nВведите команду (1/0): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			jamOn()
		case "0":
			jamOff()
		default:
			fmt.Println("Неверный ввод. Используйте 1 или 0.")
		}
	}
}

// jamOn активирует блокировку: переводит пины в OUTPUT LOW
func jamOn() {
	pin0.Output()
	pin0.Low()

	pin1.Output()
	pin1.Low()

	fmt.Println("\n[!!!] БЛОКИРОВКА АКТИВНА")
	fmt.Println("Линии Data 0 и Data 1 прижаты к GND.")
	fmt.Println("Считыватель физически не может передать данные.")
}

// jamOff деактивирует блокировку: возвращает пины в INPUT (Z-состояние)
func jamOff() {
	pin0.Input()
	pin1.Input()

	fmt.Println("\n[OK] БЛОКИРОВКА СНЯТА")
	fmt.Println("Линии свободны. Система работает в штатном режиме.")
}

// cleanup корректное завершение работы
func cleanup() {
	fmt.Println("Деактивация блокировки...")
	jamOff()
	rpio.Close()
	fmt.Println("Завершение работы.")
}
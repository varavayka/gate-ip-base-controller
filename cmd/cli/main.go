package main

import (
	"fmt"
	"strconv"
	"strings"

	"log"
	"os"

	"github.com/varavayka/gate-ip-base-controller/pkg/wiegand"
)

type CardData struct {
	FacilityCode int64
	CardNumber   int64
}
type MapHandlers struct {
}

func userInput() (string, error) {
	var userInput string
	_, err := fmt.Fscan(os.Stdin, &userInput)
	return userInput, err
}
func rawCardCodeToStruct(cardInfo string) (*CardData, error) {
	rawCard := strings.Split(cardInfo, ",")
	facilityCode, err := strconv.ParseInt(rawCard[0], 10, 0)
	if len(rawCard) < 2 {
		return nil, fmt.Errorf("[ERROR] Необходимо передавать номер карты в формате: 090,76767 без пробелов")
	}

	if err != nil {
		return nil, fmt.Errorf("Ошибка %e", err)
	}
	cardNumber, err := strconv.ParseInt(rawCard[1], 10, 0)

	if err != nil {
		return nil, fmt.Errorf("Ошибка %e", err)
	}
	var cardData = CardData{FacilityCode: facilityCode, CardNumber: cardNumber}
	return &cardData, nil
}
func cliMenu(countMenuItem int) (*int64, error) {
	fmt.Println("Опции")
	fmt.Println()

	fmt.Println("[1] Считать rfid карту")
	fmt.Println("[2] Отправить карту в контроллер")
	fmt.Println("[3] Заблокировать считыватель")

	fmt.Printf(": ")

	var options string
	fmt.Fscan(os.Stdin, &options)

	option, err := strconv.ParseInt(options, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("Необходимо выбирать только цифры")
	}
	fmt.Println(countMenuItem)
	if int(option) > countMenuItem {
		return nil, fmt.Errorf("Такого раздела не существует")
	}

	switch option {
	case 1:
		fmt.Println("Ожидаю карту...")
	case 2:
		fmt.Println("[+] Ожидаю номер карты в формате 093,12176")
		fmt.Println("- Номер карты не должен иметь пробелов")
		fmt.Println("- Номер карты не должен иметь никаких спец символов кроме [ , ]")
	case 3:
		fmt.Println("Управление:")
		fmt.Println("  1 - ЗАБЛОКИРОВАТЬ считыватель")
		fmt.Println("  0 - РАЗБЛОКИРОВАТЬ считыватель")
		fmt.Println("  Ctrl+C - Выход (авто-разблокировка)")
	}

	return &option, nil

}
func main() {
	options := map[int]func(){
		1: func() { fmt.Println(wiegand.Receiver()) },
		2: func() {
			cardString, err := userInput()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			cardData, err := rawCardCodeToStruct(cardString)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			binaryPresentCode := wiegand.EncodeWiegand(uint(cardData.FacilityCode), uint(cardData.CardNumber))
			if err := wiegand.Transmit(binaryPresentCode); err != nil {
				log.Fatal(err)
			}
		},
		3: func() { wiegand.Jammer() },
	}
	option, err := cliMenu(len(options))
	if err != nil {
		fmt.Println(err)
		return
	}

	options[int(*option)]()

}

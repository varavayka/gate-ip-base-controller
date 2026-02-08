package wiegand
import "fmt"

func EncodeWiegand(facilityCode uint, cardNumber uint) string {
		// Проверка диапазонов
	if facilityCode > 255 {
		facilityCode = 255
	}
	if cardNumber > 65535 {
		cardNumber = 65535
	}

	// Конвертируем в бинарные строки
	fcBinary := fmt.Sprintf("%08b", facilityCode)   // 8 бит
	cnBinary := fmt.Sprintf("%016b", cardNumber)    // 16 бит

	// Объединяем (без битов четности)
	payload := fcBinary + cnBinary // 24 бита

	// Вычисляем биты четности
	// Even parity для первых 12 бит (FC + первые 4 бита CN)
	evenParityBits := payload[0:12]
	evenCount := 0
	for _, bit := range evenParityBits {
		if bit == '1' {
			evenCount++
		}
	}
	evenParity := "0"
	if evenCount%2 == 1 {
		evenParity = "1"
	}

	// Odd parity для последних 12 бит (последние 12 бит CN)
	oddParityBits := payload[12:24]
	oddCount := 0
	for _, bit := range oddParityBits {
		if bit == '1' {
			oddCount++
		}
	}
	oddParity := "1"
	if oddCount%2 == 1 {
		oddParity = "0"
	}

	// Собираем полный Wiegand 26
	result := evenParity + payload + oddParity

	return result
}

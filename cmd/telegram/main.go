package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/varavayka/gate-ip-base-controller/pkg/wiegand"
	"slices"
)

type UsersMap struct {
	TelegramId string `json:"telegramId"`
	CardId     string `json:"cardId"`
}

var optionsKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📝 Считать карту"),
		tgbotapi.NewKeyboardButton("✅ Войти"),
		tgbotapi.NewKeyboardButton("🚫 Заблокировать считыватель"),
	),
	// tgbotapi.NewKeyboardButtonRow(
	// 	tgbotapi.NewKeyboardButton("4"),
	// 	tgbotapi.NewKeyboardButton("5"),
	// 	tgbotapi.NewKeyboardButton("6"),
	// ),
)
var isNotRegistredKyeboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("✅ Зарегистрироваться"),
	),
)

func contextkeyboard(description string) tgbotapi.ReplyKeyboardMarkup {
	var cardsKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(description),
		),
	)
	return cardsKeyboard
}
var userList = []int{290850674}

func main() {
	
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil { // отбрасываем если это не текстовое сообщение
			continue
		}
		userID := update.Message.From.ID
		existsUserList := slices.Contains(userList, int(userID))
		var msg tgbotapi.MessageConfig
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if !existsUserList {
			msg.Text = "Профиль не зарегистрирован!"
			msg.ReplyMarkup = isNotRegistredKyeboard

		} else {
			msg.ReplyMarkup = optionsKeyboard

		}

		switch update.Message.Text {
		case "✅ Зарегистрироваться":

			userList = append(userList, int(userID))
			msg.Text = "✅ Вы зарегистрированы!"
			msg.ReplyMarkup = optionsKeyboard
		case "📝 Считать карту":
			if existsUserList {
				msg.Text = "Ожидаю карту..."
				// Создаем канал для передачи ID карты
				cardChan := make(chan string)
				go wiegand.Receiver(cardChan)

				go func(chatID int64) {
					// Программа "зависнет" на этой строке, пока карта не считается
					cardID := <-cardChan

					// Как только получили — отправляем в Telegram
					finalMsg := tgbotapi.NewMessage(chatID, "Считана карта: "+cardID)
					bot.Send(finalMsg)
				}(update.Message.Chat.ID)
				// deleteMessage(update, bot)

			}
		case "✅ Войти":
			if existsUserList {
				msg.Text = "Укажите номер карты в формате: 090,11439"
				// wiegand.EncodeWiegand(update.Message.Text)
				
				wiegand.Transmit(update.Message.Text)

				msg.Text = "✅ Дверь открыта!"

			}
		case "🚫 Заблокировать считыватель":
			if existsUserList {
				msg.Text = "В будущем будем блокировать считыватель"

			}
		}

		if _, err := bot.Send(msg); err != nil {

			panic(err)
		}

	}

}

func deleteMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	deleteMsgConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)

	// Отправляем запрос на удаление
	_, err := bot.Request(deleteMsgConfig)
	if err != nil {
		log.Printf("Ошибка при удалении: %s", err)
	}
}
func userMapParser(pathToFile string) (*UsersMap, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, fmt.Errorf("Ошибка: %s", err)

	}
	defer file.Close()
	var data = make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		data = data[:n]
	}

	var usersMap UsersMap

	err = json.Unmarshal(data, &usersMap)
	if err != nil {
		return nil, fmt.Errorf("Ошибка: %s", err)
	}
	return &usersMap, nil
}

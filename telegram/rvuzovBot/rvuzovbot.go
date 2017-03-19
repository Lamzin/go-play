package main

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/tucnak/telebot"
	"time"

	"./db"
	"./botDB"
	"./bot"
)

var log = logging.MustGetLogger("telebot")

func main() {

	db.Init()
	botDB.Init()

	tbot, err := telebot.NewBot("318894634:AAF6k5ykqoDH_ieuCTqhkGX4fCWbxi3J2cQ")
	if err != nil {
		return
	}

	messages := make(chan telebot.Message)
	tbot.Listen(messages, 1*time.Second)

	for message := range messages {
		bot.Handle(tbot, message)
		//UniversitySuggest(tbot, message)
		//fmt.Printf("%s: %s\n", message.Chat.Username, message.Text)
	}
}

func UniversitySuggest(bot *telebot.Bot, message telebot.Message) {
	log.Info(message)
	univs, err := db.UniversitySearch(message.Text)
	if err != nil {
		log.Error(err.Error())
		return
	}
	keyboard := make([][]string, 0)
	for _, u := range univs {
		keyboard = append(keyboard, []string{fmt.Sprintf("%s, %s", u.Name, u.Abbr)})
	}

	text := fmt.Sprintf(`Список университетов по запросу "%s"`, message.Text)
	bot.SendMessage(message.Chat, text, &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			ForceReply:         true,
			Selective:          true,
			ResizeKeyboard:     true,
			HideCustomKeyboard: true,
			CustomKeyboard:     keyboard,
		},
	})
}

package botDB

import (
	"github.com/tucnak/telebot"
	"gopkg.in/mgo.v2/bson"
)

type KeyboardButtonOption struct {
	Text   string
	Action string
}

type Chat struct {
	TelegramChat telebot.Chat

	University string
	Faculty    string
	Group      string

	State    string
	Keyboard [][]KeyboardButtonOption
}

func NewChat(ID int64) (chat Chat, err error) {
	err = ctx.Chats.Find(bson.M{"telegramchat.id": ID}).One(&chat)
	ctx.Refresh("get chat", err)
	return
}

func (chat Chat) Save() error {
	_, err := ctx.Chats.Upsert(bson.M{"telegramchat.id": chat.TelegramChat.ID}, chat)
	ctx.Refresh("save chat", err)
	return err
}
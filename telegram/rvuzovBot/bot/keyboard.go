package bot

import (
	"../botDB"
)

type Keyboard [][]botDB.KeyboardButtonOption

func (k Keyboard) ToTelebotKeyboard() [][]string {
	keyboard := make([][]string, len(k))
	for i, raw := range k {
		keyboard[i] = make([]string, len(raw))
		for j, item := range raw {
			keyboard[i][j] = item.Text
		}
	}
	return keyboard
}

type ButtonItem interface {
	GetID() string
	GetName() string
}

func buildKeyboardList(items []ButtonItem) Keyboard {
	keyboard := make(Keyboard, 0)
	for _, item := range items {
		buttonOption := botDB.KeyboardButtonOption{
			Text:   item.GetName(),
			Action: item.GetID(),
		}
		keyboard = append(keyboard, []botDB.KeyboardButtonOption{buttonOption})
	}
	return keyboard
}

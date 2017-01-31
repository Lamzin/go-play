package bot

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/tucnak/telebot"

	"../botDB"
	"../db"
	"time"
)

var log = logging.MustGetLogger("bot")

type handlerCtx struct {
	Telebot *telebot.Bot

	Message telebot.Message

	chat botDB.Chat
}

func Handle(bot *telebot.Bot, m telebot.Message) error {
	log.Info(m.Text)
	ctx := handlerCtx{
		Telebot: bot,
		Message: m,
	}
	return ctx.Handle()
}

func (ctx *handlerCtx) send(text string, keyboard [][]string) error {
	return ctx.Telebot.SendMessage(ctx.Message.Chat, text, &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			ResizeKeyboard:     true,
			HideCustomKeyboard: true,
			OneTimeKeyboard:    true,
			CustomKeyboard:     keyboard,
		},
	})
}

func (ctx *handlerCtx) Handle() error {
	chat, err := botDB.NewChat(ctx.Message.Chat.ID)
	if err != nil && err.Error() != "not found" {
		ctx.Telebot.SendMessage(
			ctx.Message.Chat,
			"Произошла ошибка на сервере... Мы уже работает над этим, подождите немного и попробуйте еще раз.",
			nil,
		)
		return err
	}

	ctx.chat = chat
	ctx.chat.TelegramChat = ctx.Message.Chat
	ctx.chat.Save()

	if ctx.Message.Text == "/start" {
		return ctx.Start()
	} else if ctx.Message.Text == "/reset" {
		return ctx.Reset()
	} else if ctx.Message.Text == "/hello" {
		return ctx.Hello()
	}

	switch ctx.chat.State {
	case "universitySuggest":
		return ctx.UniversitySuggest()
	case "facultySuggest":
		return ctx.FacultySuggest()
	case "groupSuggest":
		return ctx.GroupSuggest()
	case "schedule":
		return ctx.Schedule()
	default:
		return ctx.Start()
	}
}

func (ctx *handlerCtx) Start() error {
	text := "Расписание ВУЗов - MUST HAVE приложение для студента. " +
		"Более 1 000 университетов и 200 000 довольных студентов.\n\n" +
		"Расписание можно добавить и редактировать через удобный редактор editor.rvuzov.ru :)\n\n" +
		"Больше информации смотри на официальном сайте rvuzov.ru или группе vk.com/rvuzov\n\n" +
		"Возникли проблемы - пишите нам на help@rvuzov.ru"
	ctx.send(text, nil)
	return ctx.Reset()
}

func (ctx *handlerCtx) Reset() error {
	ctx.send("Для просмотра расписания необходимо указать свой университет, факультет и группу.", nil)
	return ctx.UniversitySuggestStart()
}

func (ctx *handlerCtx) Hello() error {
	ctx.send("Hi^-^", nil)

	return nil
}

func (ctx *handlerCtx) UniversitySuggestStart() error {
	ctx.send("Введи часть названия или аббревиатуры университета", nil)
	ctx.chat.University = ""
	ctx.chat.Faculty = ""
	ctx.chat.Group = ""
	ctx.chat.State = "universitySuggest"
	ctx.chat.Save()
	return nil
}

func (ctx *handlerCtx) UniversitySuggest() error {
	for _, row := range ctx.chat.Keyboard {
		for _, item := range row {
			if item.Text == ctx.Message.Text {
				ctx.chat.State = "facultySuggest"
				ctx.chat.University = item.Action
				ctx.chat.Keyboard = nil
				ctx.chat.Save()

				ctx.send(fmt.Sprintf(`Отлично! Твой университет "%s".`, ctx.Message.Text), nil)
				return ctx.FacultySuggestStart()
			}
		}
	}

	universities, _ := db.UniversitySearch(ctx.Message.Text)
	buttons := make([]ButtonItem, len(universities))
	for i, u := range universities {
		buttons[i] = ButtonItem(u)
	}
	keyboard := buildKeyboardList(buttons)
	if len(universities) != 0 {
		ctx.send(fmt.Sprintf("Список университетов по запросу `%s`.\n\nНичего не нашёл? Измени запрос и попробуй еще раз.", ctx.Message.Text), keyboard.ToTelebotKeyboard())
	} else {
		ctx.send(fmt.Sprintf("Ничего не найдено по запросу `%s`.\n\nПопробуй изменить запрос.", ctx.Message.Text), nil)
	}
	ctx.chat.Keyboard = keyboard

	ctx.chat.Save()
	return nil
}

func (ctx *handlerCtx) FacultySuggestStart() error {
	ctx.send("Введи часть названия факультета", nil)
	return nil
}

func (ctx *handlerCtx) FacultySuggest() error {
	for _, row := range ctx.chat.Keyboard {
		for _, item := range row {
			if item.Text == ctx.Message.Text {
				ctx.chat.State = "groupSuggest"
				ctx.chat.Faculty = item.Action
				ctx.chat.Keyboard = nil
				ctx.chat.Save()

				ctx.send(fmt.Sprintf(`Отлично! Твой факультет "%s".`, ctx.Message.Text), nil)
				return ctx.GroupSuggestStart()
			}
		}
	}

	faculties, _ := db.FacultySearch(ctx.chat.University, ctx.Message.Text)

	buttons := make([]ButtonItem, len(faculties))
	for i, u := range faculties {
		buttons[i] = ButtonItem(u)
	}
	keyboard := buildKeyboardList(buttons)
	if len(faculties) != 0 {
		ctx.send(fmt.Sprintf("Список факультетов по запросу `%s`.\n\nНичего не нашёл? Измени запрос и попробуй еще раз.", ctx.Message.Text), keyboard.ToTelebotKeyboard())
	} else {
		ctx.send(fmt.Sprintf("Ничего не найдено по запросу `%s`.\n\nПопробуй изменить запрос.", ctx.Message.Text), nil)
	}
	ctx.chat.Keyboard = keyboard
	ctx.chat.Save()

	return nil
}

func (ctx *handlerCtx) GroupSuggestStart() error {
	ctx.send("Введи часть названия группы", nil)
	return nil
}

func (ctx *handlerCtx) GroupSuggest() error {
	for _, row := range ctx.chat.Keyboard {
		for _, item := range row {
			if item.Text == ctx.Message.Text {
				ctx.chat.State = "schedule"
				ctx.chat.Group = item.Action
				ctx.chat.Keyboard = nil
				ctx.chat.Save()

				ctx.send(fmt.Sprintf(`Отлично! Твоя группа "%s".`, ctx.Message.Text), nil)
				return ctx.Schedule()
			}
		}
	}

	groups, _ := db.GroupSearch(ctx.chat.Faculty, ctx.Message.Text)

	buttons := make([]ButtonItem, len(groups))
	for i, u := range groups {
		buttons[i] = ButtonItem(u)
	}
	keyboard := buildKeyboardList(buttons)
	if len(groups) != 0 {
		ctx.send(fmt.Sprintf("Список групп по запросу `%s`. Ничего не нашёл?\n\nИзмени запрос и попробуй еще раз.", ctx.Message.Text), keyboard.ToTelebotKeyboard())
	} else {
		ctx.send(fmt.Sprintf("Ничего не найдено по запросу `%s`.\n\nПопробуй изменить запрос.", ctx.Message.Text), nil)
	}
	ctx.chat.Keyboard = keyboard
	ctx.chat.Save()

	return nil
}

func (ctx *handlerCtx) Schedule() error {
	var keyboard Keyboard = Keyboard{
		[]botDB.KeyboardButtonOption{
			botDB.KeyboardButtonOption{
				Text:   "Вчера",
				Action: getHumanDate(time.Now().Add(-24 * time.Hour)),
			},
			botDB.KeyboardButtonOption{
				Text:   "Сегодня",
				Action: getHumanDate(time.Now()),
			},
			botDB.KeyboardButtonOption{
				Text:   "Завтра",
				Action: getHumanDate(time.Now().Add(24 * time.Hour)),
			},
		},
	}
	ctx.chat.Keyboard = keyboard
	ctx.chat.Save()

	for _, row := range ctx.chat.Keyboard {
		for _, item := range row {
			if item.Text == ctx.Message.Text {
				lessons, _ := db.GetLessonsByDate(ctx.chat.Group, item.Action)
				if len(lessons) == 0 {
					ctx.send(fmt.Sprintf("%s\nНету пар", item.Action), keyboard.ToTelebotKeyboard())
				}
				for _, lesson := range lessons {
					text := fmt.Sprintf("%s-%s %s\n", lesson.Time.Start, lesson.Time.End, lesson.Subject)
					for _, teacher := range lesson.Teachers {
						text += teacher.Name + "\n"
					}
					ctx.send(text, keyboard.ToTelebotKeyboard())
				}
				return nil
			}
		}
	}

	ctx.send("попробуй еще раз :)", keyboard.ToTelebotKeyboard())
	ctx.send("👨‍🎨", keyboard.ToTelebotKeyboard())


	return nil
}

func getHumanDate(t time.Time) string {
	return t.Format("02.01.2006")
}

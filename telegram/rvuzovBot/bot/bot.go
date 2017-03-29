package bot

import (
	"fmt"

	"github.com/op/go-logging"
	"github.com/tucnak/telebot"

	"time"

	"../botDB"
	"../db"
	"gopkg.in/kyokomi/emoji.v1"
	"strings"
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
	} else if ctx.Message.Text == "/today" {
		return ctx.Today()
	} else if ctx.Message.Text == "/tomorrow" {
		return ctx.Tomorrow()
	} else if ctx.Message.Text == "/help" {
		return ctx.Help()
	}

	switch ctx.chat.State {
	case "universitySuggest":
		return ctx.UniversitySuggest()
	case "facultySuggest":
		return ctx.FacultySuggest()
	case "groupSuggest":
		return ctx.GroupSuggest()
	case "schedule":
		return ctx.Today()
	default:
		return ctx.Start()
	}
}

func (ctx *handlerCtx) Start() error {
	text := `Привет! Я бот расписание – скорее всего, у меня есть твое расписание.
Я умею показывать расписание занятий на текущий день. 
Для этого мне нужно знать в каком университете, факультете и группе ты учишься. 
	
Поддерживаются следущие команды:
/start - начать диалог со мной
/reset - сбросить свои настройки
/schedule - показать расписание на сегодня`
	ctx.send(text, nil)
	return ctx.Reset()
}

func (ctx *handlerCtx) Reset() error {
	ctx.send("Для просмотра расписания необходимо указать свой университет, факультет и группу.", nil)
	return ctx.UniversitySuggestStart()
}

func (ctx *handlerCtx) UniversitySuggestStart() error {
	universities, _ := db.UniversityList()

	buttons := make([]ButtonItem, len(universities))
	for i, u := range universities {
		buttons[i] = ButtonItem(u)
	}
	keyboard := buildKeyboardList(buttons)
	ctx.chat.Keyboard = keyboard

	ctx.chat.University = ""
	ctx.chat.Faculty = ""
	ctx.chat.Group = ""
	ctx.chat.State = "universitySuggest"
	ctx.chat.Save()
	err := ctx.send("Введи часть названия или аббревиатуры университета", keyboard.ToTelebotKeyboard())
	if err != nil {
		log.Error(err.Error())
	}
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
		ctx.send(emoji.Sprintf(
			"Кажется, мы не нашли твой университет.\nТы можешь добавить его через наш "+
				"онлайн-редактор editor.rvuzov.ru или там же заказать расписание – мы его добавим :slightly_smiling_face:"), nil)
		ctx.send("Кстати, ты можешь самостоятельно добавить расписание прямо в мобильном приложение:", nil)
		ctx.send("itunes.apple.com/ru/app/raspisanie-vuzov/id631171099?mt=8", nil)
		ctx.send("play.google.com/store/apps/details?id=com.raspisaniyevuzov.app", nil)
	}
	ctx.chat.Keyboard = keyboard

	ctx.chat.Save()
	return nil
}

func (ctx *handlerCtx) FacultySuggestStart() error {
	faculties, _ := db.FacultySearch(ctx.chat.University, "")

	buttons := make([]ButtonItem, len(faculties))
	for i, u := range faculties {
		buttons[i] = ButtonItem(u)
	}
	keyboard := buildKeyboardList(buttons)
	ctx.chat.Keyboard = keyboard
	ctx.chat.Save()

	ctx.send("Введи часть названия факультета или выбери из списка", keyboard.ToTelebotKeyboard())
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
		ctx.send(emoji.Sprintf(
			"Кажется, мы не нашли твой факультет.\nТы можешь добавить его через наш "+
				"онлайн-редактор editor.rvuzov.ru или там же заказать расписание – мы его добавим :slightly_smiling_face:"), nil)
		ctx.send("Кстати, ты можешь самостоятельно добавить расписание прямо в мобильном приложение:", nil)
		ctx.send("itunes.apple.com/ru/app/raspisanie-vuzov/id631171099?mt=8", nil)
		ctx.send("play.google.com/store/apps/details?id=com.raspisaniyevuzov.app", nil)
	}
	ctx.chat.Keyboard = keyboard
	ctx.chat.Save()

	return nil
}

func (ctx *handlerCtx) GroupSuggestStart() error {
	groups, _ := db.GroupSearch(ctx.chat.Faculty, "")

	buttons := make([]ButtonItem, len(groups))
	for i, u := range groups {
		buttons[i] = ButtonItem(u)
	}
	keyboard := buildKeyboardList(buttons)
	ctx.chat.Keyboard = keyboard
	ctx.chat.Save()

	ctx.send("Введи часть названия группы", keyboard.ToTelebotKeyboard())
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
				return ctx.Today()
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
		ctx.send(emoji.Sprintf(
			"Кажется, мы не нашли твою группу.\nТы можешь добавить её через наш "+
				"онлайн-редактор editor.rvuzov.ru или там же заказать расписание – мы его добавим :slightly_smiling_face:"), nil)
		ctx.send("Кстати, ты можешь самостоятельно добавить расписание прямо в мобильном приложение:", nil)
		ctx.send("itunes.apple.com/ru/app/raspisanie-vuzov/id631171099?mt=8", nil)
		ctx.send("play.google.com/store/apps/details?id=com.raspisaniyevuzov.app", nil)
	}
	ctx.chat.Keyboard = keyboard
	ctx.chat.Save()

	return nil
}

func (ctx *handlerCtx) Today() error {
	if ctx.chat.Group == "" {
		return ctx.Start()
	}

	lessons, _ := db.GetLessonsByDate(ctx.chat.Group, getHumanDate(time.Now()))
	if len(lessons) == 0 {
		ctx.send("Сегодня можно отдохнуть! :)", nil)
	}
	for _, lesson := range lessons {
		text := fmt.Sprintf(":clock3:%s-%s\n:book:%s\n", lesson.Time.Start, lesson.Time.End, lesson.Subject)

		if len(lesson.Teachers) > 0 {
			teachers := make([]string, 0)
			for _, teacher := range lesson.Teachers {
				teachers = append(teachers, teacher.Name)
			}
			text += ":man:" + strings.Join(teachers, ", ") + "\n"
		}

		if len(lesson.Audiences) > 0 {
			audiences := make([]string, 0)
			for _, audience := range lesson.Audiences {
				audiences = append(audiences, audience.Name)
			}
			text += ":earth_americas:" + strings.Join(audiences, ", ") + "\n"
		}

		text = emoji.Sprint(text)
		ctx.send(text, nil)
	}
	return nil
}

func (ctx *handlerCtx) Tomorrow() error {
	ctx.send("Для просмотра расписания на завтра, создания заданий и других функций загружайте наши приложения:", nil)
	ctx.send("play.google.com/store/apps/details?id=com.raspisaniyevuzov.app", nil)
	ctx.send("itunes.apple.com/ru/app/raspisanie-vuzov/id631171099?mt=8", nil)
	ctx.chat.Keyboard = nil
	ctx.chat.Save()
	return nil
}

func (ctx *handlerCtx) Help() error {
	text := `
	Напиши нам в поддержку @rvuzov`
	ctx.send(text, nil)
	ctx.chat.Keyboard = nil
	ctx.chat.Save()
	return nil
}

func getHumanDate(t time.Time) string {
	return t.Format("02.01.2006")
}

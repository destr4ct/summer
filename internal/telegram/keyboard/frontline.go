package keyboard

import (
	"context"
	tga "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slog"
)

// Клавиатура (основная)
// 1. О боте == /help
// 		- Отправляет справку о боте
// 2. Обновить
// 		- Запрашивает свежие данные из бд
// 		- Отправляет summary или сообщение "нет ничего нового на момент %01"

// 3. Ресурсы
//		- Устанавливает новую клавиатуру (ресурсы)

const aboutMessage = `
Summer - бот, который умеет собирать данные из разных источников и создавать на основе них краткие сводки. Под капотом используется openai как поставщик суммаризации текса. 
`

var MainKeyboard = tga.NewReplyKeyboard(
	tga.NewKeyboardButtonRow(
		tga.NewKeyboardButton("О боте"),
		tga.NewKeyboardButton("Обновить"),
	),
	tga.NewKeyboardButtonRow(
		tga.NewKeyboardButton("Ресурсы"),
	),
)

var MainMap = HandlerMap{
	"О боте":   AboutBotHandler,
	"Обновить": UpdateHandler,
	"Ресурсы":  ResourcesHandler,
}

func AboutBotHandler(ctx context.Context) error {
	args, err := ExtractArgs(ctx)
	if err != nil {
		return err
	}

	args.Logger.Info(
		"AboutBotHandler action",
		slog.String("user", args.Username()),
	)

	message := tga.NewMessage(args.Event.Message.Chat.ID, aboutMessage)
	_, err = args.Api.Send(message)
	return err
}

func UpdateHandler(ctx context.Context) error {
	return nil
}

func ResourcesHandler(ctx context.Context) error {
	return nil
}

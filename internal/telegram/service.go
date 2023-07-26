package telegram

import (
	"context"
	"destr4ct/summer/internal/storage"
	"destr4ct/summer/internal/telegram/keyboard"
	"destr4ct/summer/pkg/utils"
	"errors"
	tga "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slog"
)

type BotService struct {
	bot              *tga.BotAPI
	keyboardSelected map[string]bool

	initialKeyboard  *tga.ReplyKeyboardMarkup
	keyboardHandlers keyboard.HandlerMap

	logger  *slog.Logger
	storage storage.SummerStorage
}

// Run обрабатывает поступающие команды до тех пор, пока
func (bs *BotService) Run(ctx context.Context, cfg tga.UpdateConfig) error {
	updates := bs.bot.GetUpdatesChan(cfg)

	for {
		select {
		case newEvent := <-updates:
			// Обрабатываем событие, логируем ошибку
			if err := bs.route(ctx, &newEvent); err != nil {
				bs.logger.Error("failed to handle event", err)
			}

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
		}
	}
}

func (bs *BotService) SetInitialKeyboard(kb tga.ReplyKeyboardMarkup) {
	bs.initialKeyboard = &kb
}

func (bs *BotService) LoadHandlers(hm keyboard.HandlerMap) {
	bs.keyboardHandlers = utils.MergeMaps(hm, bs.keyboardHandlers)
}

func (bs *BotService) route(ctx context.Context, event *tga.Update) error {
	// Пропускаем события, не связанные с сообщениями
	if event.Message == nil {
		return nil
	}

	username := event.SentFrom().UserName
	message := tga.NewMessage(event.Message.Chat.ID, "Заданное действие не определено")

	// Определяем базовое действие - запуск корневой клавиатуры
	if !(bs.initialKeyboard == nil || bs.keyboardSelected[username]) {
		if event.Message.IsCommand() && event.Message.Command() == "start" {
			bs.keyboardSelected[username] = true

			message.Text = "Доступные действия отображены на клавиатуре"
			message.ReplyMarkup = keyboard.MainKeyboard
		} else {
			message.Text = "Введите /start для начала работы"
		}

	} else {
		if action, found := bs.keyboardHandlers[event.Message.Text]; found {
			return action(keyboard.WrapContext(ctx, &keyboard.HandlerArgs{
				Api:    bs.bot,
				Logger: bs.logger,
				Event:  event,
			}))
		}

	}

	_, err := bs.bot.Send(message)
	return err
}

func GetService(token string, logger *slog.Logger, storage storage.SummerStorage) (*BotService, error) {
	hdl, err := tga.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	srv := &BotService{
		bot:              hdl,
		logger:           logger,
		keyboardHandlers: make(keyboard.HandlerMap),
		keyboardSelected: make(map[string]bool),
	}

	return srv, nil
}

func DefaultConfig() tga.UpdateConfig {
	cfg := tga.NewUpdate(0)
	cfg.Timeout = 60

	return cfg
}

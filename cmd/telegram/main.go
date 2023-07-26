package main

import (
	"context"
	"destr4ct/summer/internal/config"
	"destr4ct/summer/internal/storage/postgres"
	"destr4ct/summer/internal/telegram"
	"destr4ct/summer/internal/telegram/keyboard"
	"destr4ct/summer/pkg/logging"
)

// todo: добавить graceful shutdown
func main() {
	// Инициализируем логер, получаем конфиг
	cfg := config.Load()

	logger := logging.GetLogger(cfg.Env)
	logger.Info("Initialized the logger")

	// Открываем подключение к storage
	storage, err := postgres.GetStorage(&cfg.DBConfig)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	// Получаем сервис...
	logger.Info("Initializing the telegram service")
	botService, err := telegram.GetService(cfg.TelegramConfig.APIKey, logger, storage)

	if err != nil {
		logger.Error("failed to initialize service", err)
		return
	}

	// ... И настраиваем его
	botService.SetInitialKeyboard(keyboard.MainKeyboard)
	botService.LoadHandlers(keyboard.MainMap)

	ctx := context.Background()
	if err := botService.Run(ctx, telegram.DefaultConfig()); err != nil {
		logger.Error("failed to initialize service", err)
	}
}

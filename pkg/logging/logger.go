package logging

import (
	"golang.org/x/exp/slog"
	"os"
	"sync"
)

var (
	loggerInstance *slog.Logger
	once           = new(sync.Once)
)

func GetLogger(env string) *slog.Logger {
	once.Do(func() {

		// Возможно я добавлю больше типов в дальнейшем
		// пока есть только dev и prod
		switch env {
		case "dev":
			opts := slog.HandlerOptions{
				AddSource: true,
				Level:     slog.NewAtomicLevel(slog.DebugLevel),
			}

			jh := opts.NewJSONHandler(os.Stdout)
			loggerInstance = slog.New(jh)

		case "prod":
			opts := slog.HandlerOptions{
				Level: slog.NewAtomicLevel(slog.WarnLevel),
			}

			ch := opts.NewTextHandler(os.Stdout)
			loggerInstance = slog.New(ch)

		}

	})

	return loggerInstance
}

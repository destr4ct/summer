package keyboard

import (
	"context"
	"destr4ct/summer/internal/storage"
	"errors"
	tga "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slog"
)

var (
	argsKey   = "ctx_args"
	ErrNoArgs = errors.New("no args provided")
)

type Handler = func(context.Context) error
type HandlerMap = map[string]Handler

type HandlerArgs struct {
	Api     *tga.BotAPI
	Event   *tga.Update
	Logger  *slog.Logger
	storage storage.SummerStorage
}

func (args *HandlerArgs) Username() string {
	user := args.Event.SentFrom()
	if user == nil {
		return ""
	}
	return user.UserName
}

func ExtractArgs(ctx context.Context) (*HandlerArgs, error) {
	if args, ok := ctx.Value(argsKey).(*HandlerArgs); ok {
		return args, nil
	}
	return nil, ErrNoArgs
}

func WrapContext(ctx context.Context, args *HandlerArgs) context.Context {
	return context.WithValue(ctx, argsKey, args)
}

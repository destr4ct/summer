package queue

import "context"

const SummerKey = "q_summer"

type MessageBroker[v any] interface {
	SendMessage(context.Context, string, Message[any]) error
	GetMessages(context.Context, string) ([]*Message[v], error)
	Close() error
}

type Message[V any] struct {
	Message string
	Other   V
}

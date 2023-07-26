package rmq

import (
	"context"
	"destr4ct/summer/internal/config"
	"destr4ct/summer/internal/queue"
	"destr4ct/summer/pkg/utils"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

var con *amqp.Connection

// RabbitBroker - имплементация брокера для RabbitMQ.
// При попытке отправить/получить сообщения будет открываться канал;
// Работает только со значениями типа string;
type RabbitBroker struct {
	con *amqp.Connection
	ch  *amqp.Channel
}

// renewChannel открывает новый канал, если прошлый по каким-то причинам не был открыт/закрылся
func (rb *RabbitBroker) renewChannel() (err error) {
	if rb.ch == nil || rb.ch.IsClosed() {
		rb.ch, err = rb.con.Channel()
	}
	return
}

func (rb *RabbitBroker) queueByKey(key string) (amqp.Queue, error) {
	if err := rb.renewChannel(); err != nil {
		return amqp.Queue{}, err
	}
	return rb.ch.QueueDeclare(key, false, false, false, false, nil)
}

// SendMessage просто отправляет в заданную очередь наше сообщение. Никаких хитростей
func (rb *RabbitBroker) SendMessage(ctx context.Context, key string, msg queue.Message[any]) error {
	q, err := rb.queueByKey(key)
	if err != nil {
		return err
	}

	return rb.ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			Timestamp:   time.Now(),
			Body:        []byte(msg.Message),
			ContentType: "text/plain",
		},
	)
}

// GetMessages получает сообщения до тех пор, пока не будет сигнала от контекста на прекращение чтения
func (rb *RabbitBroker) GetMessages(ctx context.Context, key string) ([]*queue.Message[time.Time], error) {
	messages := make([]*queue.Message[time.Time], 0, 8)

	q, err := rb.queueByKey(key)
	if err != nil {
		return messages, err
	}

	readChannel, err := rb.ch.Consume(
		q.Name, "",
		true, false,
		false, false, nil,
	)
	if err != nil {
		return nil, err
	}

	for {
		select {
		case newDelivery := <-readChannel:
			if newDelivery.ContentType == "text/plain" {
				newMessage := &queue.Message[time.Time]{
					Message: string(newDelivery.Body),
					Other:   newDelivery.Timestamp,
				}
				messages = append(messages, newMessage)
			}
		case <-ctx.Done():
			return messages, nil
		}
	}
}

func (rb *RabbitBroker) Close() error {
	_ = rb.ch.Close()
	return rb.con.Close()
}

// GetBroker конструктор брокера для RabbitMQ. Обновляет подключение, если это требуется.
// Использует utils.DoWithAttempts для инициализации подключения, поскольку в некоторых ситуациях rabbitmq
// отвечает нам не сразу
func GetBroker(conf *config.Config) (*RabbitBroker, error) {
	broker := &RabbitBroker{
		con: con,
	}

	if broker.con == nil || broker.con.IsClosed() {

		datasource := fmt.Sprintf(
			"amqp://%s:%s@%s:%d",
			conf.BrokerConfig.Username,
			conf.BrokerConfig.Password,
			conf.BrokerConfig.Host,
			conf.BrokerConfig.Port,
		)

		// Делаем 5 попыток на подключение
		nc, err := utils.DoWithAttempts(func() (*amqp.Connection, error) {
			return amqp.Dial(datasource)
		}, 5)

		if err != nil {
			return nil, err
		}

		broker.con = nc
	}

	return broker, nil
}

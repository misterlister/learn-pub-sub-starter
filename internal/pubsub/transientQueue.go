package pubsub

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType routing.SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	connChan, err := conn.Channel()

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	newQueue, err := connChan.QueueDeclare(
		queueName,
		queueType == routing.Durable,
		queueType == routing.Transient,
		queueType == routing.Transient,
		false,
		nil,
	)

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = connChan.QueueBind(
		queueName,
		key,
		exchange,
		false,
		nil,
	)

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return connChan, newQueue, nil
}

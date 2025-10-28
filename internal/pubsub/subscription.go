package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType routing.SimpleQueueType,
	handler func(T),
) error {
	connChan, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)

	if err != nil {
		fmt.Println(err)
	}

	deliveryChan, err := connChan.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveryChan {
			var msg T
			err = json.Unmarshal(delivery.Body, &msg)

			if err != nil {
				fmt.Printf("error: %v", err)
				continue
			}
			handler(msg)
			delivery.Ack(false)
		}
	}()

	return nil

}

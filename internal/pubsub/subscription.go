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
	handler func(T) AckType,
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
			marshalErr := json.Unmarshal(delivery.Body, &msg)

			if marshalErr != nil {
				fmt.Printf("error: unable to unmarshal data - %v\n", marshalErr)
				delivery.Nack(false, false)
			}

			ackStatus := handler(msg)

			switch ackStatus {
			case Ack:
				fmt.Println("Ack")
				delivery.Ack(false)
			case NackRequeue:
				fmt.Println("NackRequeue")
				delivery.Nack(false, true)
			case NackDiscard:
				fmt.Println("NackDiscard")
				delivery.Nack(false, false)
			default:
				fmt.Printf("Unknown AckType=%v\n", ackStatus)
				delivery.Nack(false, false)
			}
		}
	}()

	return nil

}

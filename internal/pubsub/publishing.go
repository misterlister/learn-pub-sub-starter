package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonBytes, err := json.Marshal(val)

	if err != nil {
		return err
	}

	publishStruct := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonBytes,
	}

	ch.PublishWithContext(context.Background(), exchange, key, false, false, publishStruct)

	return nil
}

package main

import (
	"fmt"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	defer conn.Close()

	fmt.Println("Connection successful!")

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		routing.Durable,
	)

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	connChan, err := conn.Channel()

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	gamelogic.PrintServerHelp()

	running := true

	for running {

		input := gamelogic.GetInput()

		if len(input) == 0 {
			continue
		}

		switch strings.ToLower(input[0]) {
		case "pause":
			fmt.Println("Sending pause command...")

			err = pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

			if err != nil {
				fmt.Printf("error: %v", err)
				continue
			}
			fmt.Println("Game paused.")

		case "resume":
			fmt.Println("Sending resume command...")

			err = pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})

			if err != nil {
				fmt.Printf("error: %v", err)
				continue
			}
			fmt.Println("Game resumed.")

		case "quit":
			fmt.Println("Exiting game.")
			running = false

		default:
			fmt.Println("Unknown command entered.")
		}
	}
}

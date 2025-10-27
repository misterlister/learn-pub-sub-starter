package main

import (
	"fmt"
	"os"
	"os/signal"
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

	connChan, err := conn.Channel()

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	gamelogic.PrintServerHelp()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	running := true

	for running {
		select {
		case <-signalChan:
			running = false
		default:
			// no interrupt signal, continue running
		}

		input := gamelogic.GetInput()

		switch {
		case len(input) == 0:
			continue

		case strings.ToLower(input[0]) == "pause":
			fmt.Println("Sending pause command...")

			err = pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

			if err != nil {
				fmt.Printf("error: %v", err)
				continue
			}
			fmt.Println("Game paused.")

		case strings.ToLower(input[0]) == "resume":
			fmt.Println("Sending resume command...")

			err = pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})

			if err != nil {
				fmt.Printf("error: %v", err)
				continue
			}
			fmt.Println("Game resumed.")

		case strings.ToLower(input[0]) == "quit":
			fmt.Println("Exiting game...")

		default:
			fmt.Println("Unknown command entered...")
		}
	}

	fmt.Println("\nShutting down connection.")
}

package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	defer conn.Close()

	fmt.Println("Connection successful!")

	username, err := gamelogic.ClientWelcome()

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		routing.Transient,
		handlerPause(gameState),
	)

	if err != nil {
		fmt.Printf("error: could not connect to game state queue. - %v\n", err)
		return
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		routing.Durable,
		handlerWar(gameState),
	)

	if err != nil {
		fmt.Printf("error: could not connect to war declarations queue. - %v\n", err)
		return
	}

	connChan, err := conn.Channel()

	if err != nil {
		fmt.Printf("error: could not create connection channel. - %v\n", err)
		return
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+username,
		routing.ArmyMovesPrefix+".*",
		routing.Transient,
		handlerMove(gameState, connChan),
	)

	if err != nil {
		fmt.Printf("error: could not connect to move queue. - %v\n", err)
		return
	}

	running := true

	for running {

		input := gamelogic.GetInput()

		if len(input) == 0 {
			continue
		}

		switch input[0] {

		case "spawn":
			err = gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			currentMove, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = pubsub.PublishJSON(
				connChan,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+username,
				currentMove,
			)

			if err != nil {
				fmt.Printf("error: could not publish move. - %v\n", err)
				continue
			}
			fmt.Println("Move published.")

		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			running = false
		default:
			fmt.Println("Unknown command entered.")
		}
	}
}

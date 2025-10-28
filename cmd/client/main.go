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

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		routing.Transient,
	)

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
		fmt.Printf("error: %v", err)
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
			_, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
			}
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

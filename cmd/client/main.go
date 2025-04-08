package main

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/pub-sub-rabbitMQ/internal/gamelogic"
	"github.com/agustin-carnevale/pub-sub-rabbitMQ/internal/pubsub"
	"github.com/agustin-carnevale/pub-sub-rabbitMQ/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connectionStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatalln("error connection to rabbitMQ server")
	}
	defer conn.Close()
	fmt.Println("Successfully connected to RabbitMQ server..")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalln("error welcoming user")
	}

	queueName := routing.PauseKey + "." + username
	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, 0)

	gameState := gamelogic.NewGameState(username)

	for {

		inputs := gamelogic.GetInput()

		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case "spawn":
			err := gameState.CommandSpawn(inputs)
			if err != nil {
				fmt.Println("couldn't spawn unit:", err)
			}
		case "move":
			_, err := gameState.CommandMove(inputs)
			if err != nil {
				fmt.Println("couldn't move unit:", err)
			}
			fmt.Println("Successful move!")
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Sorry I couldn't process that instruction.. try again!")
		}
	}

	// wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan

	// fmt.Println("Peril client is shutting down..")
	// conn.Close()
}

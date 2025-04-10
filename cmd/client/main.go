package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

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

	// queueName := routing.PauseKey + "." + username
	// pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, 0)

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalln("error creating channel")
	}

	gameState := gamelogic.NewGameState(username)
	pauseQueueName := routing.PauseKey + "." + username
	pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, pauseQueueName, routing.PauseKey, 0, handlerPause(gameState))
	movesQueueName := routing.ArmyMovesPrefix + "." + username
	movesKey := routing.ArmyMovesPrefix + ".*"
	pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, movesQueueName, movesKey, 0, handlerMove(gameState, publishCh))
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		1,
		handlerWar(gameState, publishCh),
	)
	if err != nil {
		log.Fatalf("could not subscribe to war declarations: %v", err)
	}

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
			move, err := gameState.CommandMove(inputs)
			if err != nil {
				fmt.Println("couldn't move unit:", err)
			} else {
				pubsub.PublishJSON(publishCh, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+username, move)
				fmt.Println("Successful move!")
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			// fmt.Println("Spamming not allowed yet!")
			if len(inputs) != 2 {
				fmt.Println("Invalid number of inputs, usage: spam <n>")
				continue
			}
			n, err := strconv.Atoi(inputs[1])
			if err != nil {
				fmt.Println("Invalid input, usage: spam <n> where n is an integer")
				continue
			}
			for i := 0; i < n; i++ {
				msg := gamelogic.GetMaliciousLog()
				pubsub.PublishGob(publishCh, routing.ExchangePerilTopic, routing.GameLogSlug+"."+username, routing.GameLog{
					CurrentTime: time.Now(),
					Username:    username,
					Message:     msg,
				})
			}
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

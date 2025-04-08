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
	fmt.Println("Starting Peril server...")
	connectionStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatalln("error connection to rabbitMQ server")
	}
	defer conn.Close()
	fmt.Println("Successfully connected to RabbitMQ server.")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("error creating rabbit channel")
	}

	// pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
	// 	IsPaused: true,
	// })
	pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, "game_logs", routing.GameLogSlug, 1)

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			fmt.Println("Sending 'Pause' event..")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
		case "resume":
			fmt.Println("Sending 'Resume' event..")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
		case "quit":
			fmt.Println("Peril is shutting down..")
			return
		default:
			fmt.Println("Sorry I couldn't process that instruction.. try again!")
		}

	}

	// wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	// fmt.Println("Peril is shutting down..")
	// conn.Close()
}

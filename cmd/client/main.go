package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Peril client is shutting down..")
	conn.Close()
}

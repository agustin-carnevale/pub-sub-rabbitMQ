package main

import (
	"fmt"

	"github.com/agustin-carnevale/pub-sub-rabbitMQ/internal/gamelogic"
	"github.com/agustin-carnevale/pub-sub-rabbitMQ/internal/pubsub"
	"github.com/agustin-carnevale/pub-sub-rabbitMQ/internal/routing"
)

func handlerLog() func(gl routing.GameLog) pubsub.AckType {
	return func(gl routing.GameLog) pubsub.AckType {
		defer fmt.Print("> ")
		gamelogic.WriteLog(gl)
		return pubsub.Ack
	}
}

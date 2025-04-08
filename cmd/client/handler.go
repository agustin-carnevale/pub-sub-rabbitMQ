package main

import (
	"fmt"

	"github.com/agustin-carnevale/pub-sub-rabbitMQ/internal/gamelogic"
	"github.com/agustin-carnevale/pub-sub-rabbitMQ/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(move gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(move)
	}
}

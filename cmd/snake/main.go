package main

import (
	"log"

	snakeEngine "github.com/Rican7/gogames/gameengine/snake"
	"github.com/Rican7/gogames/internal/snake"
)

func main() {
	gameEngine := snakeEngine.NewGameEngine(20, 20)
	gameFrontEnd := snake.NewFrontEnd(gameEngine, 200, 220)
	err := gameFrontEnd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

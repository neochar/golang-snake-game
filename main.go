package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
)

func main() {
	seed := int64(1)
	//seed := time.Now().UnixNano()

	rand.New(rand.NewSource(seed))

	game := &Game{
		snake:    NewSnake(),
		food:     NewFood(),
		gameOver: false,
		delay:    startDelay,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalln(err)
	}
}

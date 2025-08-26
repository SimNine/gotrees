package main

import (
	"log"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	allDims := util.Dims{X: 640, Y: 480}
	game := game.NewGame(allDims)
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

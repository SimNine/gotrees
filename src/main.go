package main

import (
	"log"

	"github.com/SimNine/go-urfutils/src/geom"
	"github.com/SimNine/gotrees/src/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	allDims := geom.Dims[int]{X: 640, Y: 480}
	game := game.NewGame(allDims)
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

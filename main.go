package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Dims Pos
type Pos struct {
	x int
	y int
}

type Game struct {
	wordsPos    Pos
	xIncreasing bool
	yIncreasing bool

	windowSize Dims
}

func (g *Game) Init() {
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowSize(g.windowSize.x, g.windowSize.y)
}

func (g *Game) Update() error {
	if g.xIncreasing {
		g.wordsPos.x += 1
	} else {
		g.wordsPos.x -= 1
	}
	if g.yIncreasing {
		g.wordsPos.y += 1
	} else {
		g.wordsPos.y -= 1
	}

	if g.wordsPos.x > g.windowSize.x {
		g.xIncreasing = false
	} else if g.wordsPos.x <= 0 {
		g.xIncreasing = true
	}
	if g.wordsPos.y > g.windowSize.y {
		g.yIncreasing = false
	} else if g.wordsPos.y <= 0 {
		g.yIncreasing = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "ayoooo", g.wordsPos.x, g.wordsPos.y)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {

	game := &Game{
		wordsPos:    Pos{0, 0},
		xIncreasing: true,
		yIncreasing: true,
		windowSize:  Dims{320, 240},
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

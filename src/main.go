package main

import (
	"log"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/environment"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	windowSize       util.Dims
	windowRenderDims util.Dims

	environment *environment.Environment
}

func (g *Game) Init() {
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowSize(g.windowSize.X, g.windowSize.Y)
}

func (g *Game) Update() error {
	// Update the game board with any non-interactive logic
	g.environment.Update()

	// Handle mouse input
	pos := util.MakePosFromTuple(ebiten.CursorPosition())
	g.environment.SetCursorPos(pos)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.environment.MouseDown()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.environment.MouseUp()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "ayoooo", 0, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowRenderDims.X, g.windowRenderDims.Y
}

func main() {
	game := &Game{
		windowSize:       util.Dims{X: 640, Y: 480},
		windowRenderDims: util.Dims{X: 640, Y: 480},
		environment:      &environment.Environment{},
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

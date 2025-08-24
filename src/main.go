package main

import (
	"log"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	windowSize       util.Dims
	windowRenderDims util.Dims

	simulation *simulation.Simulation
}

func (g *Game) Init() {
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowSize(g.windowSize.X, g.windowSize.Y)
}

func (g *Game) Update() error {
	// Update the game board with any non-interactive logic
	g.simulation.Update()

	// Handle mouse input
	pos := util.MakePosFromTuple(ebiten.CursorPosition())
	g.simulation.SetCursorPos(pos)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.simulation.MouseDown()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.simulation.MouseUp()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.simulation.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowRenderDims.X, g.windowRenderDims.Y
}

func main() {
	allDims := util.Dims{X: 640, Y: 480}
	game := &Game{
		windowSize:       allDims,
		windowRenderDims: allDims,
		simulation: simulation.NewSimulation(
			allDims,
		),
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/localutil"
	"github.com/SimNine/gotrees/src/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	windowSize       util.Dims
	windowRenderDims util.Dims

	viewport   localutil.Viewport
	simulation *simulation.Simulation
}

func (g *Game) Init() {
	ebiten.SetWindowTitle("GeneTrees")
	ebiten.SetWindowSize(g.windowSize.X, g.windowSize.Y)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}

func (g *Game) Update() error {
	// Update the game board with any non-interactive logic
	g.simulation.Update()

	// Get mouse position and adjust for viewport
	pos := util.MakePosFromTuple(ebiten.CursorPosition())
	pos.X += g.viewport.Pos.X
	pos.Y += g.viewport.Pos.Y

	// Handle mouse input
	g.simulation.SetCursorPos(pos)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.simulation.MouseDown()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.simulation.MouseUp()
	}

	// Handle keyboard input for viewport movement
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.viewport.Pos.Y -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.viewport.Pos.Y += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.viewport.Pos.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.viewport.Pos.X += 5
	}

	// Check if the window size has changed
	w, h := ebiten.WindowSize()
	if w != g.windowSize.X || h != g.windowSize.Y {
		g.windowSize.X = w
		g.windowSize.Y = h
		g.windowRenderDims = g.windowSize
		g.viewport.Dims = g.windowRenderDims
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.simulation.Draw(screen, g.viewport)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowRenderDims.X, g.windowRenderDims.Y
}

func main() {
	allDims := util.Dims{X: 640, Y: 480}
	game := &Game{
		windowSize:       allDims,
		windowRenderDims: allDims,
		viewport: localutil.Viewport{
			Pos:  util.Pos[int]{X: 0, Y: 0},
			Dims: allDims,
		},
		simulation: simulation.NewSimulation(
			util.Dims{X: 2000, Y: 1000},
		),
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

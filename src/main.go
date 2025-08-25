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
	viewport         util.Pos[int] // Top-left corner of the viewport in world coordinates

	simulation *simulation.Simulation
}

func (g *Game) Init() {
	ebiten.SetWindowTitle("GeneTrees")
	ebiten.SetWindowSize(g.windowSize.X, g.windowSize.Y)
}

func (g *Game) Update() error {
	// Update the game board with any non-interactive logic
	g.simulation.Update()

	// Get mouse position and adjust for viewport
	pos := util.MakePosFromTuple(ebiten.CursorPosition())
	pos.X += g.viewport.X
	pos.Y += g.viewport.Y

	// Handle mouse input
	g.simulation.SetCursorPos(pos)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.simulation.MouseDown()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.simulation.MouseUp()
	}

	// Handle keyboard input for viewport movement
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.viewport.Y -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.viewport.Y += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.viewport.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.viewport.X += 5
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
		viewport:         util.Pos[int]{X: 0, Y: 0},
		simulation: simulation.NewSimulation(
			allDims,
		),
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

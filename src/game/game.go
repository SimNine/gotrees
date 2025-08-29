package game

import (
	"log"

	"github.com/SimNine/go-urfutils/src/geom"
	"github.com/SimNine/gotrees/src/simulation"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func NewGame(dims geom.Dims[int]) *Game {
	return &Game{
		windowSize:       dims,
		windowRenderDims: dims,
		viewport: geom.Viewport[int]{
			Bounds: geom.MakeBoundsFromPosAndDims(
				geom.Pos[int]{X: 0, Y: 0},
				dims,
			),
			Debug: true,
		},
		cursorWindowPos: geom.Pos[int]{X: 0, Y: 0},
		cursorPressed:   false,
		simulation: simulation.NewSimulation(
			geom.Dims[int]{X: 4000, Y: 2000},
		),
	}
}

type Game struct {
	windowSize       geom.Dims[int]
	windowRenderDims geom.Dims[int]
	viewport         geom.Viewport[int]
	cursorWindowPos  geom.Pos[int]
	prevCursorPos    geom.Pos[int]
	cursorPressed    bool

	simulation *simulation.Simulation
}

func (g *Game) Init() {
	ebiten.SetWindowTitle("GeneTrees")
	ebiten.SetWindowSize(g.windowSize.X, g.windowSize.Y)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(500)
}

func (g *Game) Update() error {
	// Update the simulation
	g.simulation.Update()

	// Handle mouse input
	g.setCursorPos(geom.MakePosFromTuple(ebiten.CursorPosition()))
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.mouseDown()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.mouseUp()
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

	// Toggle debug mode
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.viewport.Debug = !g.viewport.Debug
	}

	// Move the viewport if the cursor was dragged
	if g.cursorPressed {
		g.viewport.Pos = g.viewport.Pos.TranslatePos(g.prevCursorPos.Sub(g.cursorWindowPos))
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

func (g *Game) setCursorPos(pos geom.Pos[int]) {
	g.prevCursorPos = g.cursorWindowPos
	g.cursorWindowPos = pos
}

func (g *Game) mouseDown() {
	g.cursorPressed = true
	log.Println("Mouse down at game pos", g.viewport.ScreenToGame(g.cursorWindowPos))
}

func (g *Game) mouseUp() {
	g.cursorPressed = false
	log.Println("Mouse up at game pos", g.viewport.ScreenToGame(g.cursorWindowPos))
}

package simulation

import (
	"math/rand"
	"strconv"

	"github.com/SimNine/go-urfutils/src/geom"
	"github.com/SimNine/gotrees/src/environment"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var TICKS_PER_GENERATION = 1000

func NewSimulation(dims geom.Dims[int]) *Simulation {
	random := rand.New(rand.NewSource(0))
	return &Simulation{
		env: environment.NewEnvironment(
			random,
			dims,
		),
		random: random,
	}
}

type Simulation struct {
	random *rand.Rand

	env *environment.Environment

	tickNum       int
	generationNum int
}

func (s *Simulation) Draw(
	screen *ebiten.Image,
	viewport geom.Viewport[int],
) {
	s.env.Draw(screen, viewport)

	// Print out debug info
	printRoot := 10
	ebitenutil.DebugPrintAt(screen, "FPS: "+strconv.FormatFloat(ebiten.ActualFPS(), 'f', 3, 64), 10, printRoot)
	ebitenutil.DebugPrintAt(screen, "TPS: "+strconv.FormatFloat(ebiten.ActualTPS(), 'f', 3, 64), 10, printRoot+15)
	ebitenutil.DebugPrintAt(screen, "Generation: "+strconv.Itoa(s.generationNum), 10, printRoot+30)
	ebitenutil.DebugPrintAt(screen, "Num trees: "+strconv.Itoa(s.env.NumTrees()), 10, printRoot+45)
}

func (s *Simulation) Update() {
	s.env.Update()

	s.tickNum++
	if s.tickNum >= TICKS_PER_GENERATION {
		s.tickNum = 0
		s.generationNum++
		s.env.AdvanceGeneration()
	}
}

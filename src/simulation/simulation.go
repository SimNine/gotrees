package simulation

import (
	"math/rand"

	"github.com/SimNine/go-urfutils/src/geom"
	"github.com/SimNine/gotrees/src/environment"
	"github.com/SimNine/gotrees/src/localutil"
	"github.com/hajimehoshi/ebiten/v2"
)

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
}

func (s *Simulation) Draw(
	screen *ebiten.Image,
	viewport localutil.Viewport[int],
) {
	s.env.Draw(screen, viewport)
}

func (s *Simulation) Update() {
	s.env.Update()
}

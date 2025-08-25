package simulation

import (
	"log"
	"math/rand"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/environment"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewSimulation(dims util.Dims) *Simulation {
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

	cursorPos util.Pos[int]

	env *environment.Environment
}

func (s *Simulation) Draw(
	screen *ebiten.Image,
	viewport util.Pos[int],
) {
	s.env.Draw(screen, viewport)
}

func (s *Simulation) Update() {
}

func (s *Simulation) SetCursorPos(pos util.Pos[int]) {
	s.cursorPos = pos
}

func (s *Simulation) MouseDown() {
	log.Println("Mouse down at", s.cursorPos)
}

func (s *Simulation) MouseUp() {
	log.Println("Mouse up at", s.cursorPos)
}

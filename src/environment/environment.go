package environment

import (
	"image/color"
	"math/rand"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/environment/genetree"
	"github.com/hajimehoshi/ebiten/v2"
)

var COLOR_SKYBLUE = color.RGBA{
	R: 100,
	G: 181,
	B: 246,
	A: 255,
}

func NewEnvironment(
	random *rand.Rand,
	dims util.Dims,
) *Environment {
	return &Environment{
		random: random,

		dims: dims,

		trees: []*genetree.GeneTree{},
		sun:   []*ParticleSun{},
		rain:  []*ParticleRain{},
		// seeds: []*ParticleSeed{},

		landscape: NewLandscape(
			random,
			dims,
			200,
		),
	}
}

type Environment struct {
	random *rand.Rand

	dims util.Dims

	trees []*genetree.GeneTree
	sun   []*ParticleSun
	rain  []*ParticleRain
	// seeds []*ParticleSeed

	landscape *Landscape
}

func (e *Environment) Draw(
	screen *ebiten.Image,
	viewport util.Pos[int],
) {
	// Fill the background with blue
	screen.Fill(COLOR_SKYBLUE)

	// Draw the landscape
	e.landscape.Draw(screen, viewport)

	// Draw the trees
	for _, tree := range e.trees {
		tree.Draw(screen)
	}

	// Draw the particles
	for _, s := range e.sun {
		s.Draw(screen)
	}
	for _, r := range e.rain {
		r.Draw(screen)
	}
	// for _, s := range e.seeds {
	// 	s.Draw(screen)
	// }
}

func (e *Environment) Update() {
}

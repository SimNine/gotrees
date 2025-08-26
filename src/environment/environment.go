package environment

import (
	"image/color"
	"math/rand"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/environment/genetree"
	"github.com/SimNine/gotrees/src/localutil"
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

	env := &Environment{
		random: random,

		dims: dims,

		trees: []*genetree.GeneTree{},
		sun:   []*ParticleSun{},
		rain:  []*ParticleRain{},
		// seeds: []*ParticleSeed{},

		landscape: NewLandscape(
			random,
			dims,
			600,
		),
	}

	// Add some trees
	for i := 0; i < 10; i++ {
		xPos := random.Intn(dims.X)
		yPos := env.landscape.groundLevels[xPos]
		env.trees = append(
			env.trees,
			genetree.NewGeneTree(
				random,
				util.Pos[int]{X: xPos, Y: yPos},
			),
		)
	}

	return env
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
	viewport localutil.Viewport,
) {
	// Fill the background with blue
	screen.Fill(COLOR_SKYBLUE)

	// Draw the landscape
	e.landscape.Draw(screen, viewport)

	// Draw the trees
	for _, tree := range e.trees {
		tree.Draw(screen, viewport)
	}

	// Draw the particles
	for _, s := range e.sun {
		s.Draw(screen, viewport)
	}
	for _, r := range e.rain {
		r.Draw(screen, viewport)
	}
	// for _, s := range e.seeds {
	// 	s.Draw(screen)
	// }
}

func (e *Environment) Update() {
}

package environment

import (
	"image/color"
	"log"
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

		trees: map[*genetree.GeneTree]struct{}{},
		sun:   map[*ParticleSun]struct{}{},
		rain:  map[*ParticleRain]struct{}{},
		// seeds: map[*ParticleSeed]struct{}{},

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
		env.trees[genetree.NewGeneTree(
			random,
			util.Pos[int]{X: xPos, Y: yPos},
		)] = struct{}{}
	}

	return env
}

type Environment struct {
	random *rand.Rand

	dims util.Dims

	trees map[*genetree.GeneTree]struct{}
	sun   map[*ParticleSun]struct{}
	rain  map[*ParticleRain]struct{}
	// seeds map[*ParticleSeed]struct{}

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
	for tree := range e.trees {
		tree.Draw(screen, viewport)
	}

	// Draw the particles
	for s := range e.sun {
		s.Draw(screen, viewport)
	}
	for r := range e.rain {
		r.Draw(screen, viewport)
	}
	// for _, s := range e.seeds {
	// 	s.Draw(screen)
	// }
}

func (e *Environment) Update() {
	tps := ebiten.ActualTPS()
	log.Println("TPS:", tps)

	e.addNewSun()
	e.addNewRain()

	for sun := range e.sun {
		sun.tick()
	}
	for rain := range e.rain {
		rain.tick()
	}
	// for _, seed := range e.seeds {
	// 	seed.tick()
	// }
	e.collideSunWithGround(e.sun)
	e.collideRainWithGround(e.rain)

	log.Println("Num Sun:", len(e.sun), "Num Rain:", len(e.rain))
}

func (e *Environment) addNewSun() {
	for i := 0; i < 2; i++ {
		pct := e.random.Float32()
		pct = pct * pct
		xPos := min(int(pct*float32(e.dims.X)), e.dims.X-1)
		e.sun[NewParticleSun(
			util.Pos[int]{X: xPos, Y: 0},
		)] = struct{}{}
	}
}

func (e *Environment) addNewRain() {
	for i := 0; i < 2; i++ {
		pct := e.random.Float32()
		pct = 1.0 - pct*pct
		xPos := min(int(pct*float32(e.dims.X)), e.dims.X-1)
		e.rain[NewParticleRain(
			util.Pos[int]{X: xPos, Y: 0},
		)] = struct{}{}
	}
}

func (e *Environment) collideSunWithGround(particles map[*ParticleSun]struct{}) {
	remParticles := []*ParticleSun{}
	for p := range particles {
		if (*p).collidesWithGround(e.landscape) {
			(*p).consume()
			remParticles = append(remParticles, p)
		}
	}
	for _, p := range remParticles {
		delete(particles, p)
	}
}

func (e *Environment) collideRainWithGround(particles map[*ParticleRain]struct{}) {
	remParticles := []*ParticleRain{}
	for p := range particles {
		if (*p).collidesWithGround(e.landscape) {
			(*p).consume()
			remParticles = append(remParticles, p)
		}
	}
	for _, p := range remParticles {
		delete(particles, p)
	}
}

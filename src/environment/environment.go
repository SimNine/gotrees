package environment

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/SimNine/go-urfutils/src/geom"
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
	dims geom.Dims[int],
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
	for i := 0; i < 50; i++ {
		xPos := random.Intn(dims.X)
		yPos := env.landscape.groundLevels[xPos]
		env.trees[genetree.NewGeneTree(
			random,
			geom.Pos[int]{X: xPos, Y: yPos},
		)] = struct{}{}
	}

	return env
}

type Environment struct {
	random *rand.Rand

	dims geom.Dims[int]

	trees map[*genetree.GeneTree]struct{}
	sun   map[*ParticleSun]struct{}
	rain  map[*ParticleRain]struct{}
	// seeds map[*ParticleSeed]struct{}

	landscape *Landscape
}

func (e *Environment) Draw(
	screen *ebiten.Image,
	viewport geom.Viewport[int],
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

	// Do all stuff with particles
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
	// e.collideSunWithTrees(e.sun)
	// e.collideRainWithTrees(e.rain)

	// // Update all trees
	// for tree := range e.trees {
	// 	tree.Update()
	// }

	// Compute fitness of all trees

	// After a certain number of ticks, reproduce or kill each tree based on fitness

	log.Println("Num Sun:", len(e.sun), "Num Rain:", len(e.rain))
}

func (e *Environment) addNewSun() {
	for i := 0; i < 2; i++ {
		pct := e.random.Float32()
		pct = pct * pct
		xPos := min(int(pct*float32(e.dims.X)), e.dims.X-1)
		e.sun[NewParticleSun(
			geom.Pos[int]{X: xPos, Y: 0},
		)] = struct{}{}
	}
}

func (e *Environment) addNewRain() {
	for i := 0; i < 2; i++ {
		pct := e.random.Float32()
		pct = 1.0 - pct*pct
		xPos := min(int(pct*float32(e.dims.X)), e.dims.X-1)
		e.rain[NewParticleRain(
			geom.Pos[int]{X: xPos, Y: 0},
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
		if p.collidesWithGround(e.landscape) {
			p.consume()
			remParticles = append(remParticles, p)
		}
	}
	for _, p := range remParticles {
		delete(particles, p)
	}
}

func (e *Environment) collideSunWithTrees(particles map[*ParticleSun]struct{}) {
	remParticles := []*ParticleSun{}
	for p := range particles {
		for tree := range e.trees {
			if p.collidesWithTree(tree) {
				p.consume()
				remParticles = append(remParticles, p)
				break
			}
		}
	}
	for _, p := range remParticles {
		delete(particles, p)
	}
}

func (e *Environment) collideRainWithTrees(particles map[*ParticleRain]struct{}) {
	remParticles := []*ParticleRain{}
	for p := range particles {
		for tree := range e.trees {
			if p.collidesWithTree(tree) {
				p.consume()
				remParticles = append(remParticles, p)
				break
			}
		}
	}
	for _, p := range remParticles {
		delete(particles, p)
	}
}

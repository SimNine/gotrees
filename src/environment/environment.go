package environment

import (
	"image/color"
	"math/rand"

	"github.com/SimNine/go-urfutils/src/geom"
	"github.com/SimNine/gotrees/src/environment/genetree"
	"github.com/hajimehoshi/ebiten/v2"
)

var MINIMUM_NUM_TREES = 100
var FITNESS_ROOT_NUTRIENT_COLLECTION_PER_SIZE = 3.0
var FITNESS_ROOT_NUTRIENT_COLLECTION_PER_DEPTH = 0.03
var FITNESS_STRUCT_DECAY_PER_SIZE = 1.0
var FITNESS_ACTIVE_NODE_DECAY_PER_SIZE = 1.0
var FITNESS_INACTIVE_NODE_DECAY_PER_SIZE = 10.0
var FITNESS_REQUIREMENT_PER_CHILD_PER_NODE = 5000.0
var BASE_MUTATION_CHANCE = 0.35

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

		bounds: geom.MakeBoundsFromPosAndDims(
			geom.Pos[int]{X: 0, Y: 0},
			dims,
		),

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
	env.addNewTrees(MINIMUM_NUM_TREES)

	return env
}

type Environment struct {
	random *rand.Rand

	bounds geom.Bounds[int]

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
	e.collideSunWithTrees(e.sun)
	e.collideRainWithTrees(e.rain)

	// Check each tree's nodes for fitness changes
	for tree := range e.trees {
		allNodes := tree.GetAllNodes()
		for node := range allNodes {
			// Collide each tree's root nodes with the ground
			if node.NodeType == genetree.TREENODE_ROOT && e.bounds.Contains(node.Pos) {
				groundLevel := e.landscape.groundLevels[node.Pos.X]
				distUnderground := float64(node.Pos.Y - groundLevel)
				tree.Nutrients += FITNESS_ROOT_NUTRIENT_COLLECTION_PER_SIZE * node.Diameter * distUnderground * FITNESS_ROOT_NUTRIENT_COLLECTION_PER_DEPTH
			}

			// Remove fitness based on node size
			if node.NodeType == genetree.TREENODE_STRUCT {
				tree.Fitness -= node.Diameter * FITNESS_STRUCT_DECAY_PER_SIZE
				// } else if node.Activated {
				// 	tree.Fitness -= node.Diameter * FITNESS_ACTIVE_NODE_DECAY_PER_SIZE
			} else {
				tree.Fitness -= node.Diameter * FITNESS_ACTIVE_NODE_DECAY_PER_SIZE
			}
		}
	}

	// Update all trees
	for tree := range e.trees {
		tree.Update()
	}
}

func (e *Environment) AdvanceGeneration() {
	// For each tree
	nextGenTrees := map[*genetree.GeneTree]struct{}{}
	for tree := range e.trees {
		// If it has fitness
		if tree.Fitness > 0 {
			// Keep it
			nextGenTrees[tree] = struct{}{}

			// Try to reproduce it
			for i := tree.Fitness; i > 0; i -= float64(len(tree.GetAllNodes())) * FITNESS_REQUIREMENT_PER_CHILD_PER_NODE {
				newXPos := tree.GetRootPos().X + e.random.Intn(200) - 200
				if newXPos < 0 || newXPos >= e.bounds.Dims.X {
					continue
				}
				newPos := geom.Pos[int]{
					X: newXPos,
					Y: e.landscape.groundLevels[newXPos],
				}
				mutate := e.random.Float64() < BASE_MUTATION_CHANCE
				newTree := tree.Clone(
					newPos,
					mutate,
				)
				nextGenTrees[newTree] = struct{}{}
			}
		}
	}

	// Replace the trees with the next generation
	e.trees = nextGenTrees

	// Reset all trees
	for tree := range e.trees {
		tree.Reset()
	}

	// If there are too few trees, add some more
	if len(e.trees) < MINIMUM_NUM_TREES {
		e.addNewTrees(MINIMUM_NUM_TREES - len(e.trees))
	}
}

func (e *Environment) GetTrees() map[*genetree.GeneTree]struct{} {
	return e.trees
}

func (e *Environment) addNewTrees(num int) {
	for i := 0; i < num; i++ {
		xPos := e.random.Intn(e.bounds.Dims.X)
		yPos := e.landscape.groundLevels[xPos]
		e.addNewTree(geom.Pos[int]{X: xPos, Y: yPos})
	}
}

func (e *Environment) addNewTree(pos geom.Pos[int]) {
	e.trees[genetree.NewGeneTree(
		e.random,
		pos,
	)] = struct{}{}
}

func (e *Environment) addNewSun() {
	for i := 0; i < 2; i++ {
		pct := e.random.Float32()
		pct = pct * pct
		xPos := min(int(pct*float32(e.bounds.Dims.X)), e.bounds.Dims.X-1)
		e.sun[NewParticleSun(
			geom.Pos[int]{X: xPos, Y: 0},
		)] = struct{}{}
	}
}

func (e *Environment) addNewRain() {
	for i := 0; i < 2; i++ {
		pct := e.random.Float32()
		pct = 1.0 - pct*pct
		xPos := min(int(pct*float32(e.bounds.Dims.X)), e.bounds.Dims.X-1)
		e.rain[NewParticleRain(
			geom.Pos[int]{X: xPos, Y: 0},
		)] = struct{}{}
	}
}

func (e *Environment) collideSunWithGround(particles map[*ParticleSun]struct{}) {
	remParticles := []*ParticleSun{}
	for p := range particles {
		if e.landscape.tileType[p.pos.X][p.pos.Y] == tileTypeGround {
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
		if e.landscape.tileType[p.pos.X][p.pos.Y] == tileTypeGround {
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
			collides, nodeType := tree.DoesPointCollide(p.pos)
			if collides {
				remParticles = append(remParticles, p)
				if nodeType == genetree.TREENODE_LEAF {
					tree.Energy += p.power
				}
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
			collides, nodeType := tree.DoesPointCollide(p.pos)
			if collides && nodeType == genetree.TREENODE_RAINCATCHER {
				remParticles = append(remParticles, p)
				tree.Energy += p.power
				break
			}
		}
	}
	for _, p := range remParticles {
		delete(particles, p)
	}
}

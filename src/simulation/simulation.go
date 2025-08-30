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
	random := rand.New(rand.NewSource(int64(rand.Int())))
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
	Paused        bool
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
	ebitenutil.DebugPrintAt(screen, "Num trees: "+strconv.Itoa(len(s.env.GetTrees())), 10, printRoot+45)

	// Gather info on trees
	maxFitness := -100000000000000.0
	minFitness := 100000000000000.0
	totalFitness := 0.0
	maxNutrients := -100000000000000.0
	minNutrients := 100000000000000.0
	totalNutrients := 0.0
	maxEnergy := -100000000000000.0
	minEnergy := 100000000000000.0
	totalEnergy := 0.0
	for tree := range s.env.GetTrees() {
		if tree.Fitness > maxFitness {
			maxFitness = tree.Fitness
		}
		if tree.Fitness < minFitness {
			minFitness = tree.Fitness
		}
		totalFitness += tree.Fitness

		if tree.Nutrients > maxNutrients {
			maxNutrients = tree.Nutrients
		}
		if tree.Nutrients < minNutrients {
			minNutrients = tree.Nutrients
		}
		totalNutrients += tree.Nutrients

		if tree.Energy > maxEnergy {
			maxEnergy = tree.Energy
		}
		if tree.Energy < minEnergy {
			minEnergy = tree.Energy
		}
		totalEnergy += tree.Energy
	}
	ebitenutil.DebugPrintAt(screen, "Avg fitness: "+strconv.FormatFloat(totalFitness/float64(len(s.env.GetTrees())), 'f', 3, 64), 10, printRoot+60)
	ebitenutil.DebugPrintAt(screen, "Min fitness: "+strconv.FormatFloat(minFitness, 'f', 3, 64), 10, printRoot+75)
	ebitenutil.DebugPrintAt(screen, "Max fitness: "+strconv.FormatFloat(maxFitness, 'f', 3, 64), 10, printRoot+90)
	ebitenutil.DebugPrintAt(screen, "Avg nutrients: "+strconv.FormatFloat(totalNutrients/float64(len(s.env.GetTrees())), 'f', 3, 64), 10, printRoot+105)
	ebitenutil.DebugPrintAt(screen, "Min nutrients: "+strconv.FormatFloat(minNutrients, 'f', 3, 64), 10, printRoot+120)
	ebitenutil.DebugPrintAt(screen, "Max nutrients: "+strconv.FormatFloat(maxNutrients, 'f', 3, 64), 10, printRoot+135)
	ebitenutil.DebugPrintAt(screen, "Avg energy: "+strconv.FormatFloat(totalEnergy/float64(len(s.env.GetTrees())), 'f', 3, 64), 10, printRoot+150)
	ebitenutil.DebugPrintAt(screen, "Min energy: "+strconv.FormatFloat(minEnergy, 'f', 3, 64), 10, printRoot+165)
	ebitenutil.DebugPrintAt(screen, "Max energy: "+strconv.FormatFloat(maxEnergy, 'f', 3, 64), 10, printRoot+180)
}

func (s *Simulation) Update() {
	if !s.Paused {
		s.env.Update()

		s.tickNum++
		if s.tickNum >= TICKS_PER_GENERATION {
			s.tickNum = 0
			s.generationNum++
			s.env.AdvanceGeneration()
		}
	}
}

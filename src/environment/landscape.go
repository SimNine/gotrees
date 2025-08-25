package environment

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type tileType int

const (
	tileTypeGround tileType = iota
	tileTypeAir
)

func NewLandscape(
	random *rand.Rand,
	dims util.Dims,
	baseLevel int,
) *Landscape {
	// Create the landscape without any details
	tileTypes := make([][]tileType, dims.X)
	for x := 0; x < dims.X; x++ {
		tileTypes[x] = make([]tileType, dims.Y)
		for y := 0; y < dims.Y; y++ {
			tileTypes[x][y] = tileTypeAir
		}
	}
	landscape := &Landscape{
		groundBaseline: baseLevel,

		tileType:     tileTypes,
		groundLevels: make([]int, dims.X),
	}

	// Create landscape mathematical seeds
	landscape.groundFrequency = []float64{
		0.002,
		0.01,
		0.04,
		0.2,
		0.5,
	}
	landscape.groundAmplitude = []float64{
		random.Float64() * 500,
		random.Float64() * 200,
		random.Float64() * 80,
		random.Float64() * 5,
		random.Float64() * 5,
	}
	landscape.groundDisplacement = []float64{
		random.Float64() * 500,
		random.Float64() * 500,
		random.Float64() * 500,
		random.Float64() * 500,
		random.Float64() * 500,
	}
	landscape.groundDegree = min(len(landscape.groundFrequency), len(landscape.groundAmplitude), len(landscape.groundDisplacement))

	// Create the tiles
	for x := 0; x < len(landscape.tileType); x++ {
		landscape.groundLevels[x] = landscape.getAlgorithmicGroundLevel(x)
		for y := landscape.groundLevels[x]; y < dims.Y; y++ {
			landscape.tileType[x][y] = tileTypeGround
		}
	}

	// Create the image
	landscape.image = ebiten.NewImage(dims.X, dims.Y)
	for x := 0; x < dims.X; x++ {
		for y := 0; y < dims.Y; y++ {
			if y > landscape.groundLevels[x] {
				landscape.image.Set(x, y, color.RGBA{R: 183, G: 85, B: 23, A: 255}) // Brown
			}
		}
	}

	return landscape
}

type Landscape struct {
	image *ebiten.Image

	groundBaseline     int
	groundDegree       int
	groundFrequency    []float64
	groundAmplitude    []float64
	groundDisplacement []float64

	tileType     [][]tileType
	groundLevels []int
}

func (l *Landscape) Draw(
	screen *ebiten.Image,
	viewport util.Pos[int],
) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-viewport.X), float64(-viewport.Y))
	screen.DrawImage(l.image, op)
}

func (l *Landscape) getAlgorithmicGroundLevel(x int) int {
	sum := 0.0
	for i := 0; i < l.groundDegree; i++ {
		sum += math.Cos(l.groundFrequency[i]*float64(x)+l.groundDisplacement[i]) * l.groundAmplitude[i]
	}
	return int(sum)
}

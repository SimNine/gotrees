package genetree

import (
	"image/color"
	"math/rand"

	"github.com/SimNine/go-urfutils/src/geom"
	"github.com/SimNine/go-urfutils/src/gfx"
	"github.com/hajimehoshi/ebiten/v2"
)

var whiteImage *ebiten.Image
var whiteSubImage *ebiten.Image

func NewGeneTree(
	random *rand.Rand,
	pos geom.Pos[int],
) *GeneTree {
	tree := &GeneTree{
		fitness:   0,
		nutrients: 0,
		energy:    0,

		root: *NewTreeNodeBase(
			random,
			pos,
		),
		age: 0,
	}

	// Get the bounding box of the tree
	tree.bounds = tree.root.getMaxSubtreeBounds()

	return tree
}

type GeneTree struct {
	random *rand.Rand

	// Cached data; should be invalidated on any change
	debugImage *ebiten.Image
	bounds     geom.Bounds[int]

	fitness           int
	nutrients         int     // fitness component from soil
	energy            int     // fitness component from sunlight
	fitnessPercentile float32 // fitness as a percentile

	root TreeNode
	age  int
}

func (t *GeneTree) Draw(
	screen *ebiten.Image,
	viewport geom.Viewport[int],
) {
	if t.debugImage == nil {
		t.debugImage = gfx.EbitenCreateHollowRectangleImage(
			t.bounds.Dims,
			color.RGBA{
				R: 255,
				G: 0,
				B: 0,
				A: 255,
			},
		)
	}
	if viewport.Debug {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(t.bounds.Pos.X-viewport.Pos.X), float64(t.bounds.Pos.Y-viewport.Pos.Y))
		screen.DrawImage(t.debugImage, op)
	}

	t.root.Draw(screen, viewport)
}

func (t *GeneTree) IsPointInBounds(pos geom.Pos[int]) bool {
	// TODO
	return false
}

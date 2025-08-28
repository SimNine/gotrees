package genetree

import (
	"math/rand"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/localutil"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewGeneTree(
	random *rand.Rand,
	pos util.Pos[int],
) *GeneTree {
	return &GeneTree{
		fitness:   0,
		nutrients: 0,
		energy:    0,

		root: *NewTreeNodeBase(
			TREENODE_STRUCT,
			20,
			pos,
		),
		age: 0,
	}
}

type GeneTree struct {
	fitness           int
	nutrients         int     // fitness component from soil
	energy            int     // fitness component from sunlight
	fitnessPercentile float32 // fitness as a percentile

	root TreeNode
	age  int
}

func (t *GeneTree) Draw(
	screen *ebiten.Image,
	viewport localutil.Viewport,
) {
	t.root.Draw(screen, viewport)
}

func (t *GeneTree) IsPointInBounds(pos util.Pos[int]) bool {
	// TODO
	return false
}

package genetree

import (
	"github.com/SimNine/go-solitaire/src/util"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewGeneTree() *GeneTree {
	return &GeneTree{
		// fitness:   0,
		// nutrients: 0,
		// energy:    0,

		// root: TreeNode{
		// 	nodeType:  NodeTypeRoot,
		// 	size:      20,
		// 	dist:      0,
		// 	angle:     0,
		// 	pos:       [2]int{0, 0},
		// 	activated: true,
		// },
		// age: 0,
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

func (t *GeneTree) Draw(screen *ebiten.Image) {
	// t.root.Draw(screen, t)
}

func (t *GeneTree) IsPointInBounds(pos util.Pos[int]) bool {
	// TODO
	return false
}

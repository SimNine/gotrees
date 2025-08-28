package genetree

import (
	"image/color"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/localutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func NewTreeNodeBase(
	nodeType NodeType,
	diameter float32,
	pos util.Pos[int],
) *TreeNode {
	return &TreeNode{
		nodeType:  nodeType,
		diameter:  diameter,
		dist:      0,
		angle:     0,
		pos:       pos,
		activated: true,
	}
}

func NewTreeNode(
	nodeType NodeType,
	diameter float32,
	dist float32,
	angle float32,
	pos util.Pos[int],
) *TreeNode {
	return &TreeNode{
		nodeType:  nodeType,
		diameter:  diameter,
		dist:      dist,
		angle:     angle,
		pos:       pos,
		activated: true,
	}
}

type TreeNode struct {
	nodeType NodeType
	diameter float32       // diameter
	dist     float32       // distance from parent node
	angle    float32       // angle (clockwise) from directly below parent (in radians)
	pos      util.Pos[int] // position of the top-left corner

	activated bool // whether this node has been used, or is vestigial
}

func (n *TreeNode) Draw(
	screen *ebiten.Image,
	viewport localutil.Viewport,
) {
	screenPos := viewport.GameToScreen(n.pos)
	vector.DrawFilledCircle(
		screen,
		float32(screenPos.X),
		float32(screenPos.Y),
		n.diameter/2,
		NODE_COLORS[n.nodeType],
		false,
	)

	if viewport.Debug {
		vector.DrawFilledRect(
			screen,
			float32(screenPos.X)-n.diameter/2,
			float32(screenPos.Y)-n.diameter/2,
			n.diameter,
			n.diameter,
			color.RGBA{R: 255, G: 0, B: 0, A: 10},
			false,
		)
	}
}

func (n *TreeNode) IsPointInBounds(pos util.Pos[int]) bool {
	// TODO
	return false
}

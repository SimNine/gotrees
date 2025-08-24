package genetree

import "github.com/SimNine/go-solitaire/src/util"

type TreeNode struct {
	nodeType NodeType
	size     float32       // diameter
	dist     float32       // distance from parent node
	angle    float32       // angle (clockwise) from directly below parent (in radians)
	pos      util.Pos[int] // position of the top-left corner

	activated bool // whether this node has been used, or is vestigial
}

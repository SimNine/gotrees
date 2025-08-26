package genetree

import "image/color"

type NodeType int

const (
	TREENODE_STRUCT NodeType = iota
	TREENODE_LEAF
	TREENODE_ROOT
	TREENODE_RAINCATCHER
	TREENODE_SEEDDROPPER
)

var NODE_COLORS = map[NodeType]color.RGBA{
	TREENODE_STRUCT:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
	TREENODE_LEAF:        color.RGBA{R: 52, G: 237, B: 52, A: 255},
	TREENODE_ROOT:        color.RGBA{R: 137, G: 47, B: 4, A: 255},
	TREENODE_RAINCATCHER: color.RGBA{R: 0, G: 0, B: 150, A: 255},
	TREENODE_SEEDDROPPER: color.RGBA{R: 0, G: 100, B: 100, A: 255},
}

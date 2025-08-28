package genetree

import (
	"image/color"
	"math/rand"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/localutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const NODE_MIN_DIAMETER = 10.0
const NODE_MIN_DISTANCE = 40.0

const NODE_MUTATE_CHANCE_TYPE = 0.15
const NODE_MUTATE_CHANCE_DIAMETER = 0.30
const NODE_MUTATE_CHANCE_DELETE_NODE = 0.10
const NODE_MUTATE_CHANCE_ADD_NODE = 0.30
const NODE_MUTATE_CHANCE_ANGLE = 0.25
const NODE_MUTATE_CHANCE_DISTANCE = 0.15

// public double mutationTreeBaseChance = 0.3;

func NewTreeNodeBase(
	random *rand.Rand,
	pos util.Pos[int],
) *TreeNode {
	return NewTreeNode(
		random,
		nil,
		NodeType(random.Intn(len(NODE_COLORS))),
		float32((random.Float64()*9.0)+NODE_MIN_DIAMETER),
		0,
		0,
		pos,
		true,
	)
}

func NewTreeNode(
	random *rand.Rand,
	children map[*TreeNode]struct{},
	nodeType NodeType,
	diameter float32,
	dist float32,
	angle float32,
	pos util.Pos[int],
	mutate bool,
) *TreeNode {
	treeNode := &TreeNode{
		random:    random,
		children:  children,
		nodeType:  nodeType,
		diameter:  diameter,
		dist:      dist,
		angle:     angle,
		pos:       pos,
		activated: true,
	}

	if mutate {
		treeNode.mutate()
	}

	return treeNode
}

type TreeNode struct {
	random *rand.Rand

	children map[*TreeNode]struct{}

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

func (n *TreeNode) mutate() {
	// Chance of mutating the node's type
	if n.random.Float32() < NODE_MUTATE_CHANCE_TYPE {
		newType := NodeType(n.random.Intn(len(NODE_COLORS)))

		// Additional chance to re-roll if not a struct
		for newType != TREENODE_STRUCT {
			if n.random.Float32() > 0.4 {
				break
			}
			newType = NodeType(n.random.Intn(len(NODE_COLORS)))
		}

		// Remove children if type is no longer a struct
		if newType != TREENODE_STRUCT {
			n.children = map[*TreeNode]struct{}{}
		}

		// Set the new type
		n.nodeType = newType
	}

	// Chance of mutating this node's diameter
	if n.random.Float32() < NODE_MUTATE_CHANCE_DIAMETER {
		diameterChange := n.random.Float32()*16.0 - 8.0
		n.diameter += diameterChange
		if n.diameter < NODE_MIN_DIAMETER {
			n.diameter = NODE_MIN_DIAMETER
		}
	}

	// Chance to lose each child node, otherwise mutate them
	if len(n.children) > 0 {
		nodesToDel := []*TreeNode{}
		for child := range n.children {
			if n.random.Float32() < NODE_MUTATE_CHANCE_DELETE_NODE {
				// Lose this child
				nodesToDel = append(nodesToDel, child)
			} else {
				// Mutate this child
				child.mutate()
			}
		}
		for _, child := range nodesToDel {
			delete(n.children, child)
		}
	}

	// Chance to add child nodes if this is a struct
	if n.nodeType == TREENODE_STRUCT {
		for {
			if n.random.Float32() < NODE_MUTATE_CHANCE_ADD_NODE {
				// Add a new child node
				child := NewTreeNodeBase(n.random, util.Pos[int]{X: n.pos.X, Y: n.pos.Y + int(n.diameter)})
				n.children[child] = struct{}{}
				child.mutate()
			} else {
				break
			}
		}
	}

	// Chance to mutate angle between this node and its parent
	if n.random.Float32() < NODE_MUTATE_CHANCE_ANGLE {
		angleChange := (n.random.Float32() * 30) - 30.0
		n.angle += angleChange
	}

	// Chance to mutate distance between this node and its parent
	if n.random.Float32() < NODE_MUTATE_CHANCE_DISTANCE {
		distChange := (n.random.Float32() * 30.0) - 30.0
		n.dist += distChange
		if n.dist < NODE_MIN_DISTANCE {
			n.dist = NODE_MIN_DISTANCE
		}
	}
}

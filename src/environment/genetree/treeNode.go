package genetree

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/SimNine/go-urfutils/src/geom"
	"github.com/SimNine/go-urfutils/src/gfx"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const NODE_MIN_DIAMETER = 10.0
const NODE_MIN_DISTANCE = 40.0

const NODE_MUTATE_CHANCE_TYPE = 0.15
const NODE_MUTATE_CHANCE_DIAMETER = 0.30
const NODE_MUTATE_CHANCE_DELETE_NODE = 0.10
const NODE_MUTATE_CHANCE_ADD_NODE = 0.25
const NODE_MUTATE_CHANCE_ANGLE = 0.25
const NODE_MUTATE_CHANCE_DISTANCE = 0.15

// public double mutationTreeBaseChance = 0.3;

func NewTreeNodeBase(
	random *rand.Rand,
	pos geom.Pos[int],
) *TreeNode {
	treeNode := NewTreeNode(
		random,
		map[*TreeNode]struct{}{},
		TREENODE_STRUCT,
		(random.Float64()*9.0)+NODE_MIN_DIAMETER,
		0,
		0,
		pos,
		true,
	)

	// init the position of all child nodes
	for child := range treeNode.children {
		child.initPosition(pos)
	}

	return treeNode
}

func NewTreeNode(
	random *rand.Rand,
	children map[*TreeNode]struct{},
	nodeType NodeType,
	diameter float64,
	dist float64,
	angle float64,
	pos geom.Pos[int],
	mutate bool,
) *TreeNode {
	treeNode := &TreeNode{
		random:    random,
		children:  children,
		nodeType:  nodeType,
		diameter:  diameter,
		dist:      dist,
		angleRads: angle,
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

	// Cached data; should be invalidated on any change
	image      *ebiten.Image
	debugImage *ebiten.Image

	children map[*TreeNode]struct{}

	nodeType  NodeType
	diameter  float64       // diameter
	dist      float64       // distance from parent node
	angleRads float64       // angle (clockwise) from directly below parent (in radians)
	pos       geom.Pos[int] // position of the top-left corner

	activated bool // whether this node has been used, or is vestigial
}

func (n *TreeNode) Draw(
	screen *ebiten.Image,
	viewport geom.Viewport[int],
) {
	centerPos := viewport.GameToScreen(n.pos)
	topleftPos := centerPos.Sub(geom.Pos[int]{X: int(n.diameter / 2), Y: int(n.diameter / 2)})

	// Draw a line from this node to each child
	for child := range n.children {
		startPos := centerPos
		endPos := viewport.GameToScreen(geom.Pos[int]{X: child.pos.X, Y: child.pos.Y})
		vector.StrokeLine(
			screen,
			float32(startPos.X),
			float32(startPos.Y),
			float32(endPos.X),
			float32(endPos.Y),
			float32(n.diameter/10),
			color.RGBA{R: 139, G: 69, B: 19, A: 255}, // brown
			false,
		)
	}

	// Draw this node
	if n.image == nil {
		imgSize := int(math.Ceil(n.diameter))
		if imgSize < 1 {
			imgSize = 1
		}
		n.image = ebiten.NewImage(imgSize, imgSize)
		vector.DrawFilledCircle(
			n.image,
			float32(imgSize)/2,
			float32(imgSize)/2,
			float32(n.diameter)/2,
			NODE_COLORS[n.nodeType],
			false,
		)
	}
	drawOptions := &ebiten.DrawImageOptions{}
	drawOptions.GeoM.Translate(float64(topleftPos.X), float64(topleftPos.Y))
	screen.DrawImage(n.image, drawOptions)

	// Draw the debug image if in debug mode
	if viewport.Debug {
		if n.debugImage == nil {
			n.debugImage = gfx.EbitenCreateHollowRectangleImage(
				geom.Dims[int]{
					X: int(math.Ceil(n.diameter)),
					Y: int(math.Ceil(n.diameter)),
				},
				color.RGBA{R: 255, G: 255, B: 0, A: 255},
			)
		}
		drawOptions := &ebiten.DrawImageOptions{}
		drawOptions.GeoM.Translate(float64(topleftPos.X), float64(topleftPos.Y))
		screen.DrawImage(n.debugImage, drawOptions)
	}

	// Draw all child nodes
	for child := range n.children {
		child.Draw(screen, viewport)
	}
}

func (n *TreeNode) DoesPointCollideRecursive(pos geom.Pos[int]) bool {
	xDiff := math.Abs(pos.ToFloatPos().X - n.pos.ToFloatPos().X)
	yDiff := math.Abs(pos.ToFloatPos().Y - n.pos.ToFloatPos().Y)
	collides := n.diameter >= math.Sqrt(math.Pow(xDiff, 2)+math.Pow(yDiff, 2))
	if collides {
		return true
	}

	for child := range n.children {
		if child.DoesPointCollideRecursive(pos) {
			return true
		}
	}

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
		diameterChange := n.random.Float64()*16.0 - 8.0
		n.diameter += diameterChange
	}
	if n.diameter < NODE_MIN_DIAMETER {
		n.diameter = NODE_MIN_DIAMETER
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
				child := NewTreeNodeBase(n.random, geom.Pos[int]{X: n.pos.X, Y: n.pos.Y})
				n.children[child] = struct{}{}
				child.mutate()
			} else {
				break
			}
		}
	}

	// Chance to mutate angle between this node and its parent
	if n.random.Float32() < NODE_MUTATE_CHANCE_ANGLE {
		angleChange := (n.random.Float64() * 0.5) - 0.5
		n.angleRads += angleChange
	}

	// Chance to mutate distance between this node and its parent
	if n.random.Float32() < NODE_MUTATE_CHANCE_DISTANCE {
		distChange := (n.random.Float64() * 30.0) - 30.0
		n.dist += distChange
	}
	if n.dist < NODE_MIN_DISTANCE {
		n.dist = NODE_MIN_DISTANCE
	}
}

func (n *TreeNode) initPosition(parentPos geom.Pos[int]) {
	// Calculate the position of this node based on its parent
	n.pos.X = parentPos.X + int(n.dist*math.Cos(n.angleRads))
	n.pos.Y = parentPos.Y + int(n.dist*math.Sin(n.angleRads))

	// Adjust the position of all child nodes
	for child := range n.children {
		child.initPosition(n.pos)
	}
}

func (n *TreeNode) getMaxSubtreeBounds() geom.Bounds[int] {
	halfDiameter := int(math.Ceil(n.diameter / 2))
	bounds := geom.MakeBoundsFromPosAndDims(
		n.pos.Sub(geom.Pos[int]{X: halfDiameter, Y: halfDiameter}),
		geom.Dims[int]{X: int(math.Ceil(n.diameter)), Y: int(math.Ceil(n.diameter))},
	)
	for child := range n.children {
		childBounds := child.getMaxSubtreeBounds()
		bounds = bounds.Union(childBounds)
	}
	return bounds
}

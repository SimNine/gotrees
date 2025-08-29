package genetree

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/localutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var whiteImage *ebiten.Image
var whiteSubImage *ebiten.Image

func NewGeneTree(
	random *rand.Rand,
	pos util.Pos[int],
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
	topLeft, bottomRight := tree.root.getMaxSubtreeBounds()
	tree.topLeft = topLeft
	tree.bottomRight = bottomRight

	return tree
}

type GeneTree struct {
	random *rand.Rand

	// Cached data; should be invalidated on any change
	debugImage  *ebiten.Image
	topLeft     util.Pos[int]
	bottomRight util.Pos[int]

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
	if t.debugImage == nil {

		if whiteImage == nil {
			whiteImage = ebiten.NewImage(3, 3)
			whiteImage.Fill(color.White)
			whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
		}

		imgSizeX := t.bottomRight.X - t.topLeft.X
		imgSizeY := t.bottomRight.Y - t.topLeft.Y
		t.debugImage = ebiten.NewImage(
			imgSizeX,
			imgSizeY,
		)

		var path vector.Path
		path.MoveTo(0, 0)
		path.LineTo(float32(imgSizeX), 0)
		path.LineTo(float32(imgSizeX), float32(imgSizeY))
		path.LineTo(0, float32(imgSizeY))
		path.Close()

		op := &vector.StrokeOptions{}
		op.LineCap = vector.LineCapSquare
		op.LineJoin = vector.LineJoinMiter
		op.MiterLimit = 5.0
		op.Width = 2.0
		vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, op)
		for i := range vs {
			vs[i].SrcX = 1
			vs[i].SrcY = 1
			vs[i].ColorR = 1
			vs[i].ColorG = 0
			vs[i].ColorB = 0
			vs[i].ColorA = 1
		}
		t.debugImage.DrawTriangles(vs, is, whiteSubImage, nil)
	}
	if viewport.Debug {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(t.topLeft.X-viewport.Pos.X), float64(t.topLeft.Y-viewport.Pos.Y))
		screen.DrawImage(t.debugImage, op)
	}

	t.root.Draw(screen, viewport)
}

func (t *GeneTree) IsPointInBounds(pos util.Pos[int]) bool {
	// TODO
	return false
}

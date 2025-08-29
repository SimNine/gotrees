package localutil

import (
	"image"
	"image/color"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var whiteImage *ebiten.Image = nil
var whiteSubImage *ebiten.Image = nil

func CreateHollowRectangleImage(
	dims util.Dims,
	c color.Color,
) *ebiten.Image {
	if whiteImage == nil {
		whiteImage = ebiten.NewImage(3, 3)
		whiteImage.Fill(color.White)
		whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	}

	ret := ebiten.NewImage(
		dims.X,
		dims.Y,
	)

	var path vector.Path
	path.MoveTo(0, 0)
	path.LineTo(float32(dims.X), 0)
	path.LineTo(float32(dims.X), float32(dims.Y))
	path.LineTo(0, float32(dims.Y))
	path.Close()

	op := &vector.StrokeOptions{}
	op.LineCap = vector.LineCapSquare
	op.LineJoin = vector.LineJoinMiter
	op.MiterLimit = 5.0
	op.Width = 2.0
	floatColor := color.RGBAModel.Convert(c).(color.RGBA)
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, op)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(floatColor.R) / 255.0
		vs[i].ColorG = float32(floatColor.G) / 255.0
		vs[i].ColorB = float32(floatColor.B) / 255.0
		vs[i].ColorA = float32(floatColor.A) / 255.0
	}
	ret.DrawTriangles(vs, is, whiteSubImage, nil)
	return ret
}

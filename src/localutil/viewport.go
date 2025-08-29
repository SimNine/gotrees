package localutil

import (
	"github.com/SimNine/go-urfutils/src/geom"
)

type Viewport struct {
	Pos   geom.Pos[int]  // Top-left corner of the viewport in world coordinates
	Dims  geom.Dims[int] // Dimensions of the viewport in pixels
	Debug bool
}

func (v *Viewport) ScreenToGame(pos geom.Pos[int]) geom.Pos[int] {
	return pos.TranslatePos(v.Pos)
}

func (v *Viewport) GameToScreen(pos geom.Pos[int]) geom.Pos[int] {
	return pos.Sub(v.Pos)
}

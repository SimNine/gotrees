package localutil

import (
	"github.com/SimNine/go-urfutils/src/geom"
)

type Viewport[N geom.Number] struct {
	Pos   geom.Pos[N]  // Top-left corner of the viewport in world coordinates
	Dims  geom.Dims[N] // Dimensions of the viewport in pixels
	Debug bool
}

func (v *Viewport[N]) ScreenToGame(pos geom.Pos[N]) geom.Pos[N] {
	return pos.TranslatePos(v.Pos)
}

func (v *Viewport[N]) GameToScreen(pos geom.Pos[N]) geom.Pos[N] {
	return pos.Sub(v.Pos)
}

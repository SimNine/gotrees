package localutil

import urfutils "github.com/SimNine/go-urfutils/src"

type Viewport struct {
	Pos   urfutils.Pos[int] // Top-left corner of the viewport in world coordinates
	Dims  urfutils.Dims     // Dimensions of the viewport in pixels
	Debug bool
}

func (v *Viewport) ScreenToGame(pos urfutils.Pos[int]) urfutils.Pos[int] {
	return pos.TranslatePos(v.Pos)
}

func (v *Viewport) GameToScreen(pos urfutils.Pos[int]) urfutils.Pos[int] {
	return pos.Sub(v.Pos)
}

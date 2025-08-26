package localutil

import "github.com/SimNine/go-solitaire/src/util"

type Viewport struct {
	Pos  util.Pos[int] // Top-left corner of the viewport in world coordinates
	Dims util.Dims     // Dimensions of the viewport in pixels
}

func (v *Viewport) ScreenToGame(pos util.Pos[int]) util.Pos[int] {
	return pos.TranslatePos(v.Pos)
}

func (v *Viewport) GameToScreen(pos util.Pos[int]) util.Pos[int] {
	return pos.Sub(v.Pos)
}

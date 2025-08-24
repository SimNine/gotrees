package environment

import (
	"log"

	"github.com/SimNine/go-solitaire/src/util"
)

type Environment struct {
	dims util.Dims

	cursorPos util.Pos[int]
}

func (e *Environment) Update() {
}

func (e *Environment) SetCursorPos(pos util.Pos[int]) {
	e.cursorPos = pos
}

func (e *Environment) MouseDown() {
	log.Println("Mouse down at", e.cursorPos)
}

func (e *Environment) MouseUp() {
	log.Println("Mouse up at", e.cursorPos)
}

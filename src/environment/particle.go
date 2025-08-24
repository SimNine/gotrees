package environment

import "github.com/SimNine/go-solitaire/src/util"

type Particle struct {
	pos        util.Pos[int]
	power      int
	isConsumed bool
}

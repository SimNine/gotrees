package environment

import (
	"image/color"

	"github.com/SimNine/go-solitaire/src/util"
	"github.com/SimNine/gotrees/src/environment/genetree"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//////////////////
// Particle
//////////////////

type Particle interface {
	tick()
	collidesWithTree(*genetree.GeneTree) bool
	collidesWithNode(*genetree.TreeNode) bool
	consume()
}

func newParticle(pos util.Pos[int], power int, color color.RGBA) baseParticle {
	return baseParticle{
		pos:        pos,
		power:      power,
		isConsumed: false,
		color:      color,
	}
}

type baseParticle struct {
	pos        util.Pos[int]
	power      int
	isConsumed bool
	color      color.RGBA
}

func (b *baseParticle) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(b.pos.X), float32(b.pos.Y), 2, b.color, false)

	// TODO
	// if (GeneTrees.debug) {
	// 	g.drawString("" + power, x - xScr, y - yScr);
	// }
}

func (p *baseParticle) collidesWithTree(t *genetree.GeneTree) bool {
	if t == nil {
		return false
	}

	return t.IsPointInBounds(p.pos)
}

func (p *baseParticle) collidesWithNode(n *genetree.TreeNode) bool {
	if n == nil {
		return false
	}

	return n.IsPointInBounds(p.pos)
}

func (p *baseParticle) consume() {
	p.isConsumed = true
	p.color = color.RGBA{R: 0, G: 0, B: 0, A: 0}
}

//////////////////
// ParticleRain
//////////////////

func NewParticleRain(pos util.Pos[int], power int) *ParticleRain {
	return &ParticleRain{
		baseParticle: newParticle(pos, power, color.RGBA{R: 0, G: 0, B: 255, A: 255}),
	}
}

type ParticleRain struct {
	baseParticle
}

func (p *ParticleRain) tick() {
	p.pos.Y += 1
	p.power += 1 // TODO: update
}

//////////////////
// ParticleSun
//////////////////

func NewParticleSun(pos util.Pos[int], power int) *ParticleSun {
	return &ParticleSun{
		baseParticle: newParticle(pos, power, color.RGBA{R: 255, G: 255, B: 0, A: 255}),
	}
}

type ParticleSun struct {
	baseParticle
}

func (p *ParticleSun) tick() {
	p.pos.Y += 1
	p.power -= 1 // TODO: update
}

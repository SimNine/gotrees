package environment

import (
	"image/color"

	"github.com/SimNine/go-urfutils/src/geom"
	"github.com/SimNine/gotrees/src/environment/genetree"
	"github.com/SimNine/gotrees/src/localutil"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const SUN_BASE_POWER = 70000
const SUN_POWER_PER_TICK = -45
const RAIN_BASE_POWER = -50000
const RAIN_POWER_PER_TICK = 55

//////////////////
// Particle
//////////////////

type Particle interface {
	tick()
	collidesWithTree(*genetree.GeneTree) bool
	collidesWithNode(*genetree.TreeNode) bool
	collidesWithGround(*Landscape) bool
	consume()
}

func newParticle(pos geom.Pos[int], power int, color color.RGBA) baseParticle {
	return baseParticle{
		pos:        pos,
		power:      power,
		isConsumed: false,
		color:      color,
	}
}

type baseParticle struct {
	pos        geom.Pos[int]
	power      int
	isConsumed bool
	color      color.RGBA
}

func (b *baseParticle) Draw(
	screen *ebiten.Image,
	viewport localutil.Viewport,
) {
	screenPos := viewport.GameToScreen(b.pos)
	vector.DrawFilledCircle(screen, float32(screenPos.X), float32(screenPos.Y), 2, b.color, false)

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

func (p *baseParticle) collidesWithGround(l *Landscape) bool {
	if l == nil {
		return false
	}

	return l.tileType[p.pos.X][p.pos.Y] == tileTypeGround
}

func (p *baseParticle) consume() {
	p.isConsumed = true
	p.color = color.RGBA{R: 0, G: 0, B: 0, A: 0}
}

//////////////////
// ParticleRain
//////////////////

func NewParticleRain(pos geom.Pos[int]) *ParticleRain {
	return &ParticleRain{
		baseParticle: newParticle(pos, RAIN_BASE_POWER, color.RGBA{R: 0, G: 0, B: 255, A: 255}),
	}
}

type ParticleRain struct {
	baseParticle
}

func (p *ParticleRain) tick() {
	p.pos.Y += 1
	p.power += RAIN_POWER_PER_TICK
}

//////////////////
// ParticleSun
//////////////////

func NewParticleSun(pos geom.Pos[int]) *ParticleSun {
	return &ParticleSun{
		baseParticle: newParticle(pos, SUN_BASE_POWER, color.RGBA{R: 255, G: 255, B: 0, A: 255}),
	}
}

type ParticleSun struct {
	baseParticle
}

func (p *ParticleSun) tick() {
	p.pos.Y += 1
	p.power += SUN_POWER_PER_TICK
}

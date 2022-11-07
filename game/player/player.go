package player

import (
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/util"
)

type Player struct {
	// Note that these are block coordinates, not pixel coordinates
	X, Y                 float64
	xVelocity, yVelocity float64
}

type MovementVector struct {
	Left, Right, Up, Down bool
}

// FIXME: consider frame delta time in equations
func (p *Player) Update(movement MovementVector) {
	var deltaX, deltaY float64 = 0, 0

	if movement.Left {
		deltaX -= 1
	}
	if movement.Right {
		deltaX += 1
	}

	if movement.Up {
		deltaY -= 1
	}
	if movement.Down {
		deltaY += 1
	}

	p.xVelocity += deltaX * config.PlayerSpeed
	p.yVelocity += deltaY * config.PlayerSpeed

	p.X += p.xVelocity
	p.Y += p.yVelocity

	p.X = util.Clamp(p.X, 0, float64(config.WorldWidth))
	p.Y = util.Clamp(p.Y, 0, float64(config.WorldHeight))

	p.xVelocity *= 0.90
	p.yVelocity *= 0.90
}

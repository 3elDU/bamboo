package player

import (
	"math"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/util"
	"github.com/3elDU/bamboo/world"
)

func (player *Player) speed() float64 {
	return math.Max(math.Abs(player.xVelocity), math.Abs(player.yVelocity))
}

func (player *Player) updateMovementDirection() {
	if player.speed() < 0.01 {
		return
	}

	var movementSide MovementDirection

	if math.Abs(player.xVelocity) > math.Abs(player.yVelocity) {
		movementSide = Left // horizontal movement
	} else {
		movementSide = Up // vertical movement
	}

	switch movementSide {
	case Left:
		if player.xVelocity > 0 {
			player.movementDirection = Right
		} else {
			player.movementDirection = Left
		}
	case Up:
		if player.yVelocity > 0 {
			player.movementDirection = Down
		} else {
			player.movementDirection = Up
		}
	}
}

// Collision points for each block are specified in local space ( e.g. relative to the block itself ),
// so for collision to work we need to convert them to global space first
func convertToGlobalSpace(block types.Block, points [4]types.Coords2f) [4]types.Coords2f {
	return [4]types.Coords2f{
		{X: points[0].X + float64(block.Coords().X), Y: points[0].Y + float64(block.Coords().Y)},
		{X: points[1].X + float64(block.Coords().X), Y: points[1].Y + float64(block.Coords().Y)},
		{X: points[2].X + float64(block.Coords().X), Y: points[2].Y + float64(block.Coords().Y)},
		{X: points[3].X + float64(block.Coords().X), Y: points[3].Y + float64(block.Coords().Y)},
	}
}

// check collision between player and the blocks
// returns collision value for each corner
func collidePlayer(origin types.Coords2f, world types.World) (collisions [4]bool) {
	playerCollisionPoints := [4]types.Coords2f{
		{X: origin.X - .25, Y: origin.Y - .25},
		{X: origin.X + .25, Y: origin.Y - .25},
		{X: origin.X - .25, Y: origin.Y + .4},
		{X: origin.X + .25, Y: origin.Y + .4},
	}

	for i, point := range playerCollisionPoints {
		interactive, isInteractive := world.BlockAt(uint64(point.X), uint64(point.Y)).(types.InteractiveBlock)
		if isInteractive {
			interactive.Interact(world, origin)
		}

		block, isCollidable := world.BlockAt(uint64(point.X), uint64(point.Y)).(types.CollidableBlock)
		if !isCollidable {
			continue
		}
		if !block.Collidable() {
			continue
		}

		blockCollisionPoints := convertToGlobalSpace(block, block.CollisionPoints())

		var blockCollisions [4]bool
		blockCollisions[0] = point.X < blockCollisionPoints[3].X || point.Y < blockCollisionPoints[3].Y
		blockCollisions[1] = point.X > blockCollisionPoints[2].X || point.Y < blockCollisionPoints[2].Y
		blockCollisions[2] = point.X < blockCollisionPoints[1].X || point.Y > blockCollisionPoints[1].Y
		blockCollisions[3] = point.X > blockCollisionPoints[0].X || point.Y > blockCollisionPoints[0].Y
		// if current point collides with any corner of the block, set the collision to true
		collisions[i] = anyOf(blockCollisions)
	}

	return
}

// returns true if any of collisions is true
func anyOf(collisions [4]bool) bool {
	for _, collision := range collisions {
		if collision {
			return true
		}
	}
	return false
}

func countCollisions(collisions [4]bool) (count uint) {
	for _, collision := range collisions {
		if collision {
			count++
		}
	}
	return
}

// Update updates the player physics and animation
// FIXME: consider frame delta time in equations
func (player *Player) Update(movement MovementVector, world *world.World) {
	dx, dy := movement.ToFloat()

	player.xVelocity += dx * config.PlayerSpeed
	player.yVelocity += dy * config.PlayerSpeed

	// if player somehow got stuck in the block, skip collision check
	if !anyOf(collidePlayer(types.Coords2f{X: player.X, Y: player.Y}, world)) {
		// check for collisions on X axis
		if anyOf(collidePlayer(types.Coords2f{X: player.X + player.xVelocity, Y: player.Y}, world)) {
			player.xVelocity = 0
		}
		// check for collisions on Y axis
		if anyOf(collidePlayer(types.Coords2f{X: player.X, Y: player.Y + player.yVelocity}, world)) {
			player.yVelocity = 0
		}
		// check for corner collisions
		if countCollisions(collidePlayer(types.Coords2f{X: player.X + player.xVelocity, Y: player.Y + player.yVelocity}, world)) == 1 {
			// "bounce" off the corner
			player.xVelocity = -player.xVelocity * 0.1
			player.yVelocity = -player.yVelocity * 0.1
		}
	}

	// multiply velocity by block speed modifier
	speedModifier := 1.0
	if block, ok := world.BlockAt(uint64(player.X), uint64(player.Y)).(types.CollidableBlock); ok {
		speedModifier = block.PlayerSpeed()
	}

	player.X += player.xVelocity * speedModifier
	player.Y += player.yVelocity * speedModifier

	player.X = util.Clamp(player.X, 0, float64(config.WorldWidth))
	player.Y = util.Clamp(player.Y, 0, float64(config.WorldHeight))

	player.updateMovementDirection()

	player.xVelocity *= 0.75
	player.yVelocity *= 0.75
}

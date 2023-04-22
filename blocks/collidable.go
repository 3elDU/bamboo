package blocks

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(CollidableBlockState{})
}

type CollidableBlockState struct {
	Collidable      bool
	CollisionPoints [4]types.Vec2f
	PlayerSpeed     float64
}

type collidableBlock struct {
	collidable bool

	// Optional
	// Each collision point is a coordinate in world space
	collisionPoints [4]types.Vec2f

	// How fast player could move through this block
	// Calculated by basePlayerSpeed * playerSpeed
	// Applicable only if collidable is false
	playerSpeed float64
}

func (b *collidableBlock) Collidable() bool {
	return b.collidable
}

func defaultCollisionPoints() [4]types.Vec2f {
	return [4]types.Vec2f{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
	}
}

func (b *collidableBlock) CollisionPoints() [4]types.Vec2f {
	/*
		return [4]types.Vec2f{
			{X: float64(b.x) + b.collisionPoints[0].X, Y: float64(b.y) + b.collisionPoints[0].Y},
			{X: float64(b.x) + b.collisionPoints[1].X, Y: float64(b.y) + b.collisionPoints[1].Y},
			{X: float64(b.x) + b.collisionPoints[2].X, Y: float64(b.y) + b.collisionPoints[2].Y},
			{X: float64(b.x) + b.collisionPoints[3].X, Y: float64(b.y) + b.collisionPoints[3].Y},
		}
	*/
	return b.collisionPoints
}

func (b *collidableBlock) PlayerSpeed() float64 {
	return b.playerSpeed
}

func (b *collidableBlock) State() interface{} {
	return CollidableBlockState{
		Collidable:      b.collidable,
		CollisionPoints: b.collisionPoints,
		PlayerSpeed:     b.playerSpeed,
	}
}

func (b *collidableBlock) LoadState(s interface{}) {
	state := s.(CollidableBlockState)
	b.collidable = state.Collidable
	b.collisionPoints = state.CollisionPoints
	b.playerSpeed = state.PlayerSpeed
}

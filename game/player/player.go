package player

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/engine/world"
	"github.com/3elDU/bamboo/util"
	"github.com/google/uuid"
)

type Player struct {
	// Note that these are block coordinates, not pixel coordinates
	X, Y                 float64
	xVelocity, yVelocity float64
}

type MovementVector struct {
	Left, Right, Up, Down bool
}

func LoadPlayer(id uuid.UUID) *Player {
	saveDir := filepath.Join(config.WorldSaveDirectory, id.String())

	f, err := os.Open(filepath.Join(saveDir, "player.gob"))
	if err != nil {
		// if file does not exist, just return an empty object
		return &Player{X: float64(config.PlayerStartX), Y: float64(config.PlayerStartY)}
	}

	player := new(Player)
	if err := gob.NewDecoder(f).Decode(player); err != nil {
		log.Panicf("LoadPlayer() - failed to decode metadata - %v", err)
	}

	return player
}

// check collision for four angles
// if player collides, reject the movement
// if player is somehow stuck in the block, temporary disable collision for that point
func collide(origin util.Coords2f, world *world.World) (sides [4]bool) {
	playerCollisionPoints := [4]util.Coords2f{
		{X: origin.X - .25, Y: origin.Y - .25},
		{X: origin.X + .25, Y: origin.Y - .25},
		{X: origin.X - .25, Y: origin.Y + .4},
		{X: origin.X + .25, Y: origin.Y + .4},
	}

	for i, point := range playerCollisionPoints {
		block, err := world.BlockAt(uint64(point.X), uint64(point.Y))
		if err != nil {
			continue
		}

		if !block.Collidable() {
			continue
		}

		blockCollisionPoints := block.CollisionPoints()

		// collisions local for this angle
		var localPointCollisions [4]bool
		localPointCollisions[0] = point.X < blockCollisionPoints[3].X || point.Y < blockCollisionPoints[3].Y
		localPointCollisions[1] = point.X > blockCollisionPoints[2].X || point.Y < blockCollisionPoints[2].Y
		localPointCollisions[2] = point.X < blockCollisionPoints[1].X || point.Y > blockCollisionPoints[1].Y
		localPointCollisions[3] = point.X > blockCollisionPoints[0].X || point.Y > blockCollisionPoints[0].Y
		// if current angle collides with any angle of the block, set the collision to true
		sides[i] = anyOf(localPointCollisions)
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

// FIXME: consider frame delta time in equations
func (p *Player) Update(movement MovementVector, world *world.World) {
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

	var skipCollisionCheck bool

	// if player somehow got stuck in the block, temporary disable collision check
	collisions := collide(util.Coords2f{X: p.X, Y: p.Y}, world)
	if anyOf(collisions) {
		skipCollisionCheck = true
	}

	if !skipCollisionCheck {
		if anyOf(collide(util.Coords2f{X: p.X + p.xVelocity, Y: p.Y}, world)) {
			p.xVelocity = 0
		}
		// check for collisions on Y axis
		if anyOf(collide(util.Coords2f{X: p.X, Y: p.Y + p.yVelocity}, world)) {
			p.yVelocity = 0
		}
		// check for corner collisions
		if countCollisions(collide(util.Coords2f{X: p.X + p.xVelocity, Y: p.Y + p.yVelocity}, world)) == 1 {
			p.xVelocity = -p.xVelocity * 0.1
			p.yVelocity = -p.yVelocity * 0.1
		}
	}

	// multiply velocity by block speed modifier
	speedModifier := 1.0
	block, err := world.BlockAt(uint64(p.X), uint64(p.Y))
	if err == nil {
		speedModifier = block.PlayerSpeed()
	}

	p.X += p.xVelocity * speedModifier
	p.Y += p.yVelocity * speedModifier

	p.X = util.Clamp(p.X, 0, float64(config.WorldWidth))
	p.Y = util.Clamp(p.Y, 0, float64(config.WorldHeight))

	p.xVelocity *= 0.90
	p.yVelocity *= 0.90
}

func (p *Player) Save(id uuid.UUID) error {
	saveDir := filepath.Join(config.WorldSaveDirectory, id.String())

	// make a save directory, if it doesn't exist yet
	os.Mkdir(saveDir, os.ModePerm)

	// open world metadata file
	f, err := os.Create(filepath.Join(saveDir, "player.gob"))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(p); err != nil {
		return err
	}

	log.Println("Player.Save() - saved")

	return nil
}

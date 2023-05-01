package player

import (
	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/world"
	"log"
	"math/rand"
	"time"
)

type Player struct {
	// Note that these are block coordinates, not pixel coordinates
	X, Y                 float64
	xVelocity, yVelocity float64

	movementDirection MovementDirection
	animationFrame    uint8
	lastFrameChange   time.Time

	// Storing the selected world, so that we know what sub-world the player is currently in
	// Used to determine what sub-world to load
	SelectedWorld types.Save
}

type MovementDirection uint8

const (
	Left MovementDirection = iota
	Right
	Up
	Down
)

type MovementVector struct {
	Left, Right, Up, Down bool
}

func (mvec MovementVector) ToFloat() (vx, vy float64) {
	if mvec.Left {
		vx -= 1
	}
	if mvec.Right {
		vx += 1
	}

	if mvec.Up {
		vy -= 1
	}
	if mvec.Down {
		vy += 1
	}

	return
}

func isValidSpawnpoint(blockType types.BlockType) bool {
	validBlocks := []types.BlockType{
		blocks.Sand, blocks.Grass, blocks.ShortGrass, blocks.TallGrass, blocks.Flowers, blocks.RedMushroom, blocks.WhiteMushroom,
		blocks.CaveFloor,
	}

	for _, blockType2 := range validBlocks {
		if blockType == blockType2 {
			return true
		}
	}
	return false
}

// Creates a new player, picking a valid spawn point
func NewPlayer(w types.World) *Player {
	// use the same seed for reproducible spawnpoint generation
	rng := rand.New(rand.NewSource(1))

	x, y := 0, 0
	it := 1
	for {
		// Pick X and Y coordinates from 256 to 768
		x = rng.Intn(768-256) + 256
		y = rng.Intn(768-256) + 256

		// create a new chunk so that we don't overwrite chunks in the world
		c := world.NewChunk(uint64(x)/16, uint64(y)/16)
		w.Generator().GenerateImmediately(c)
		blockType := c.At(uint(x%16), uint(y%16)).Type()

		if isValidSpawnpoint(blockType) {
			break
		}

		it++
	}
	log.Printf("picked spawn point (%v, %v), took %v iterations", x, y, it)

	return &Player{X: float64(x), Y: float64(y), SelectedWorld: w.Metadata()}
}

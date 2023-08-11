package player

import (
	"log"
	"math/rand"
	"time"

	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/world"
)

type Player struct {
	// Note that these are block coordinates, not pixel coordinates
	X, Y                 float64
	xVelocity, yVelocity float64
	input                MovementVector

	MovementDirection MovementDirection
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
		types.SandBlock, types.GrassBlock, types.ShortGrassBlock, types.TallGrassBlock, types.FlowersBlock, types.RedMushroomBlock, types.WhiteMushroomBlock,
		types.CaveFloorBlock,
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
	worldWidth := int(w.Size().X)
	worldHeight := int(w.Size().Y)
	it := 1
	for {
		// Pick X and Y coordinates from 1/4 of the world to 3/4 of the world
		// E.g. if the world is 1024 blocks in size, coordinates would be in range from 256 to 768
		x = rng.Intn(worldWidth/4*3-worldWidth/4) + worldWidth/4
		y = rng.Intn(worldHeight/4*3-worldHeight/4) + worldHeight/4

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
	w.SetPlayerSpawnPoint(uint64(x), uint64(y))

	return &Player{X: float64(x), Y: float64(y), SelectedWorld: w.Metadata()}
}

package config

import (
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/world_type"
	"log"
)

// All values with type uint64 are measured in ticks, unless noted otherwise
// 1 second == 60 ticks
const (
	PerlinNoiseScaleFactor float64 = 128

	AssetDirectory = "./assets/"

	PlayerSpeed    float64 = 0.02
	PlayerInfoFile         = "player.gob"

	WorldSaveDirectory        = "./saves/"
	WorldInfoFile             = "world.gob"
	WorldAutosaveDelay uint64 = 3600
	ChunkUnloadDelay   uint64 = 600

	UIScaling float64 = 2

	OverworldSize = 1024
	Cave1Size     = 256
)

func SizeForWorldType(world world_type.WorldType) types.Vec2u {
	switch world {
	case world_type.Overworld:
		return types.Vec2u{X: OverworldSize, Y: OverworldSize}
	case world_type.Cave:
		return types.Vec2u{X: Cave1Size, Y: Cave1Size}
	}

	log.Printf("Unable to retrieve world size for world type %v", world)
	return types.Vec2u{X: 1024, Y: 1024}
}

package config

// All values with type uint64 are measured in ticks, unless noted otherwise
// 1 second == 60 ticks
const (
	AssetDirectory                 = "./assets/"
	PerlinNoiseScaleFactor float64 = 128
	PlayerSpeed            float64 = 0.02

	WorldWidth  uint64 = 1024
	WorldHeight uint64 = 1024

	PlayerStartX   = WorldWidth / 2 // initial player position, when the world is created
	PlayerStartY   = WorldHeight / 2
	PlayerInfoFile = "player.gob"

	WorldSaveDirectory        = "./saves/"
	WorldInfoFile             = "world.gob"
	WorldAutosaveDelay uint64 = 3600
	ChunkUnloadDelay   uint64 = 600

	UIScaling float64 = 3
)

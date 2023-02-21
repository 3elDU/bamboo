package config

// All values with type uint64 are measured in ticks, unless noted otherwise
// 1 second == 60 ticks
const (
	AssetDirectory         string  = "./assets/"
	PerlinNoiseScaleFactor float64 = 128
	PlayerSpeed            float64 = 0.02

	WorldWidth  uint64 = 1024
	WorldHeight uint64 = 1024

	// Initial player position, when the world is created
	PlayerStartX uint64 = WorldWidth / 2
	PlayerStartY uint64 = WorldHeight / 2

	WorldSaveDirectory string = "./saves/"
	WorldInfoFile      string = "world.gob"
	WorldAutosaveDelay uint64 = 3600
	ChunkUnloadDelay   uint64 = 600
)

var (
	FontSize float64 = 32
)

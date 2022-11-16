package config

// All values with type uint64 are measured in ticks, unless noted otherwise
// 1 second == 60 ticks
const (
	AssetDirectory         string  = "./assets/"
	PerlinNoiseScaleFactor float64 = 128
	PlayerSpeed            float64 = 0.05

	WorldWidth  int64 = 1024
	WorldHeight int64 = 1024

	// Initial player position, when the world is created
	PlayerStartX float64 = float64(WorldWidth) / 2
	PlayerStartY float64 = float64(WorldHeight) / 2

	WorldSaveDirectory string = "./saves/"
	WorldInfoFile      string = "world.gob"
	WorldAutosaveDelay uint64 = 3600
	ChunkUnloadDelay   uint64 = 600
)

var (
	FontSize float64 = 32
)

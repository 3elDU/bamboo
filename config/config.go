package config

const (
	AssetDirectory         string  = "./assets/"
	PerlinNoiseScaleFactor float64 = 128
	PlayerSpeed            float64 = 0.05

	WorldWidth  int64 = 1024
	WorldHeight int64 = 1024

	WorldSaveDirectory string = "./saves/"
	WorldInfoFile      string = "world.gob"
	WorldAutosaveDelay uint64 = 3600 // each 60 seconds
)

var (
	FontSize float64 = 32
)

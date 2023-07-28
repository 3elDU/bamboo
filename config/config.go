package config

import (
	"os"
	"path/filepath"
)

// Set externally at build time
var (
	GitCommit string = "unknown"
	GitTag    string = "unknown"

	BuildMachine string = "unknown"
	BuildDate    string = "unknown"
)

// All values with type uint64 are measured in ticks, unless noted otherwise
// 1 second == 60 ticks
const (
	PerlinNoiseScaleFactor float64 = 128

	PlayerSpeed    float64 = 0.02
	PlayerInfoFile         = "player.gob"

	WorldSaveDirectory        = "./saves/"
	WorldInfoFile             = "world.gob"
	WorldAutosaveDelay uint64 = 3600
	ChunkUnloadDelay   uint64 = 600

	InventoryFile       = "inventory.gob"
	SlotSize      uint8 = 50

	UIScaling float64 = 2

	OverworldSize = 1024
	Cave1Size     = 256
)

var AssetDirectory = "assets"

func init() {
	// Find the asset directory by looking at os.Args[0].
	// This is necessary because asset directory location is relative to the binary,
	// and the CWD could be different from the directory where the game binary is.
	AssetDirectory = filepath.Join(filepath.Dir(os.Args[0]), AssetDirectory)
}

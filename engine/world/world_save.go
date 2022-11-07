// Everything related to world saves

package world

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/game/player"
	"github.com/google/uuid"
)

func init() {
	// create saves directory, if it doesn't exist yet
	os.Mkdir(config.WorldSaveDirectory, os.ModePerm)
}

// structure with metadata, representing a world save
type WorldSave struct {
	Name string    // world name as from the user
	UUID uuid.UUID // internal unique world id, for identification purposes
	Seed int64
	Size int64 // in bytes

	// Note that player field is unused in the world itself
	// It is written on World.Save() - specifically for saving the player position
	Player player.Player
}

// Loads a world with given uuid
func Load(id uuid.UUID) (*World, player.Player) {
	saveDir := filepath.Join(config.WorldSaveDirectory, id.String())

	f, err := os.Open(filepath.Join(saveDir, "world.gob"))
	if err != nil {
		log.Panicf("World.Load() - invalid file descriptor - %v", err)
	}

	decoder := gob.NewDecoder(f)
	metadata := new(WorldSave)
	if err := decoder.Decode(metadata); err != nil {
		log.Panicf("World.Load() - failed to decode metadata - %v", err)
	}

	log.Printf("World.Load() - loaded metadata; seed - %v", metadata.Seed)

	return NewWorld(metadata.Name, metadata.UUID, metadata.Seed), metadata.Player
}

// NOTE: world folder is named after the UUID, not after the world name
// that is, to avoid folder collision
func (w *World) Save(player player.Player) error {
	saveDir := filepath.Join(config.WorldSaveDirectory, w.Metadata.UUID.String())

	// write player position
	w.Metadata.Player = player

	// make a save directory, if it doesn't exist yet
	os.Mkdir(saveDir, os.ModePerm)

	// open world metadata file
	f, err := os.Create(filepath.Join(saveDir, "world.gob"))
	if err != nil {
		return err
	}
	defer f.Close()

	// encode the metadata to it
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(w.Metadata); err != nil {
		return err
	}

	log.Println("World.Save() - saved")

	return nil
}

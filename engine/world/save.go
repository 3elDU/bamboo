// Everything related to saving/loading world

package world

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/game/player"
	"github.com/3elDU/bamboo/util"
	"github.com/google/uuid"
)

func init() {
	// create saves directory, if it doesn't exist yet
	os.Mkdir(config.WorldSaveDirectory, os.ModePerm)
}

// WorldSaverLoader maintains chunk loading/saving queue,
// while doing actual work on separate goroutine,
// so we don't have any freezes on the main thread
type WorldSaverLoader struct {
	Metadata WorldSave

	saveRequests     chan Chunk
	loadRequestsPool map[util.Coords2i]bool
	// loadRequestsPool keeps track of currently requested chunks,
	// so that one same chunk can't be requested twice
	loadRequests chan util.Coords2i
	loaded       chan *Chunk
}

func NewWorldSaverLoader(metadata WorldSave) *WorldSaverLoader {
	return &WorldSaverLoader{
		Metadata: metadata,

		saveRequests:     make(chan Chunk, 1024),
		loadRequestsPool: make(map[util.Coords2i]bool),
		loadRequests:     make(chan util.Coords2i, 256),
		loaded:           make(chan *Chunk),
	}
}

func (sl *WorldSaverLoader) runSaver() {
	for {
		chunk := <-sl.saveRequests

		// FIXME: This is probably not very save
		if err := chunk.Save(sl.Metadata.UUID); err != nil {
			log.Panicf("sl.runSaver() - error while saving chunk - %v", err)
		}
	}
}

func (sl *WorldSaverLoader) runLoader() {
	for {
		request := <-sl.loadRequests

		c := LoadChunk(sl.Metadata.UUID, request.X, request.Y)
		if c == nil {
			// if the requested chunk doesn't exist, simply ignore the error and skip it
			continue
		}
		log.Printf("WorldSaverLoader.runLoader() - loaded chunk %v; %v from disk", request.X, request.Y)

		sl.loaded <- c
	}
}

func (sl *WorldSaverLoader) Run() {
	go sl.runSaver()
	go sl.runLoader()
}

// Returns newly loaded chunk
// If there is no pending chunks, returns nil
func (sl *WorldSaverLoader) Receive() *Chunk {
	select {
	case c := <-sl.loaded:
		delete(sl.loadRequestsPool, c.Coords())
		return c
	default:
		return nil
	}
}

// Pushes chunk save request to the queue
func (sl *WorldSaverLoader) Save(chunk *Chunk) {
	sl.saveRequests <- *chunk
}

// Pushes chunk load reuqest to the queue
func (sl *WorldSaverLoader) Load(cx, cy int64) {
	coords := util.Coords2i{X: cx, Y: cy}
	if sl.loadRequestsPool[coords] {
		return
	}
	sl.loadRequestsPool[coords] = true
	sl.loadRequests <- coords
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

// block structures contain unexported fields.
// that makes it impossible to serialize them through gob
// so, we need to convert it first
// here, SavedBlock comes in handy
// it contains all required data
// + optional metadata, that can be written individually by each block
type SavedBlock struct {
	Type  BlockType
	State []byte
}

// represents chunk on the disk
// all chunks are converted to this structure before saving
type SavedChunk struct {
	X, Y int64
	Data [16][16][3]SavedBlock
}

// Loads a world with given uuid
func LoadWorld(id uuid.UUID) (*World, player.Player) {
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

	// loop over all loaded chunks, saving modified ones to the disk
	for _, chunk := range w.chunks {
		if err := chunk.Save(w.Metadata.UUID); err != nil {
			return err
		}
	}

	log.Println("World.Save() - saved")

	return nil
}

func ChunkExistsOnDisk(id uuid.UUID, x, y int64) bool {
	path := filepath.Join(config.WorldSaveDirectory, id.String(),
		fmt.Sprintf("chunk_%v_%v.gob", x, y))

	_, err := os.Stat(path)
	return err == nil
}

// if saved chunk doesn't exist, returns nil
func LoadChunk(id uuid.UUID, x, y int64) *Chunk {
	path := filepath.Join(config.WorldSaveDirectory, id.String(),
		fmt.Sprintf("chunk_%v_%v.gob", x, y))

	if _, err := os.Stat(path); err == nil {
		f, err := os.Open(path)
		if err != nil {
			log.Panicf("LoadChunk() - failed to open a file - %v", err)
		}

		savedChunk := new(SavedChunk)
		if err := gob.NewDecoder(f).Decode(savedChunk); err != nil {
			log.Panicf("LoadChunk() - failed to decode a chunk - %v", err)
		}

		c := NewChunk(x, y)

		// decode blocks
		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				for z := 0; z < 3; z++ {
					b := GetBlockByID(savedChunk.Data[x][y][z].Type)
					if err := b.LoadState(savedChunk.Data[x][y][z].State); err != nil {
						log.Panicf("LoadChunk() - Block.LoadState() failed - %v", err)
					}
					if err := c.SetBlock(x, y, Layer(z), b); err != nil {
						log.Panicf("LoadChunk() - Chunkk.SetBlock() failed - %v", err)
					}
				}
			}
		}

		// mark chunk as unmodified, to avoid recursive loading/saving
		c.modified = false
		log.Printf("LoadChunk() - loaded chunk %v; %v from disk", x, y)
		return c
	}

	return nil
}

func (c *Chunk) Save(id uuid.UUID) error {
	// if chunk wasn't modified, saving is unnecessary
	if !c.modified {
		return nil
	}

	path := filepath.Join(config.WorldSaveDirectory, id.String(),
		fmt.Sprintf("chunk_%v_%v.gob", c.x, c.y))

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// serialize the blockjs
	var blocks [16][16][3]SavedBlock
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			stack := c.blocks[x][y]
			for z, block := range []Block{stack.Bottom, stack.Ground, stack.Top} {
				blocks[x][y][z] = SavedBlock{
					Type:  block.Type(),
					State: block.State(),
				}
			}
		}
	}
	chunk := SavedChunk{
		X: c.x, Y: c.y,
		Data: blocks,
	}

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(chunk); err != nil {
		return err
	}

	c.modified = false
	log.Printf("Chunk.Save() - %v; %v", c.x, c.y)
	return nil
}

func DeleteWorld(id uuid.UUID) {
	path := filepath.Join(config.WorldSaveDirectory, id.String())
	if err := os.RemoveAll(path); err != nil {
		log.Panicf("Failed to delete world %v - %v", id, err)
	}
}

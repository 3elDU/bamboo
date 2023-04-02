// Everything related to saving/loading world

package world

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/google/uuid"
)

func init() {
	// create saves directory, if it doesn't exist yet
	os.Mkdir(config.WorldSaveDirectory, os.ModePerm)
}

// SaverLoader maintains chunk loading/saving queue,
// while doing actual work on separate goroutine,
// so we don't have any freezes on the main thread
type SaverLoader struct {
	Metadata Save

	saveRequests     chan Chunk
	loadRequestsPool map[types.Coords2u]bool
	// loadRequestsPool keeps track of currently requested chunks,
	// so that one same chunk can't be requested twice
	loadRequests chan types.Coords2u
	loaded       chan *Chunk
}

func NewWorldSaverLoader(metadata Save) *SaverLoader {
	return &SaverLoader{
		Metadata: metadata,

		saveRequests:     make(chan Chunk, 1024),
		loadRequestsPool: make(map[types.Coords2u]bool),
		loadRequests:     make(chan types.Coords2u, 256),
		loaded:           make(chan *Chunk),
	}
}

func (sl *SaverLoader) runSaver() {
	for {
		chunk := <-sl.saveRequests

		// FIXME: This is probably not very save
		chunk.Save(sl.Metadata.UUID)
	}
}

func (sl *SaverLoader) runLoader() {
	for {
		request := <-sl.loadRequests

		c := LoadChunk(sl.Metadata.UUID, request.X, request.Y)
		if c == nil {
			// if the requested chunk doesn't exist, simply ignore the error and skip it
			continue
		}
		log.Printf("SaverLoader.runLoader() - loaded chunk %v; %v from disk", request.X, request.Y)

		sl.loaded <- c
	}
}

func (sl *SaverLoader) Run() {
	go sl.runSaver()
	go sl.runLoader()
}

// Returns newly loaded chunk
// If there is no pending chunks, returns nil
func (sl *SaverLoader) Receive() *Chunk {
	select {
	case c := <-sl.loaded:
		delete(sl.loadRequestsPool, c.Coords())
		return c
	default:
		return nil
	}
}

// Pushes chunk save request to the queue
func (sl *SaverLoader) Save(chunk *Chunk) {
	sl.saveRequests <- *chunk
}

// Pushes chunk load reuqest to the queue
func (sl *SaverLoader) Load(cx, cy uint64) {
	coords := types.Coords2u{X: cx, Y: cy}
	if sl.loadRequestsPool[coords] {
		return
	}
	sl.loadRequestsPool[coords] = true
	sl.loadRequests <- coords
}

// structure with metadata, representing a world save
type Save struct {
	Name string    // world name as from the user
	UUID uuid.UUID // internal unique world id, for identification purposes
	Seed int64
	Size int64 // in bytes
}

// block structures contain unexported fields.
// that makes it impossible to serialize them through gob
// so, we need to convert it first
// here, SavedBlock comes in handy
// it contains all required data
// + optional metadata, that can be written individually by each block
type SavedBlock struct {
	Type  types.BlockType
	State interface{}
}

// represents chunk on the disk
// all chunks are converted to this structure before saving
type SavedChunk struct {
	X, Y uint64
	Data [16][16]SavedBlock
}

// Loads a world with given uuid
func LoadWorld(id uuid.UUID) *World {
	saveDir := filepath.Join(config.WorldSaveDirectory, id.String())

	f, err := os.Open(filepath.Join(saveDir, config.WorldInfoFile))
	if err != nil {
		log.Panicf("LoadWorld() - invalid file descriptor - %v", err)
	}

	decoder := gob.NewDecoder(f)
	metadata := new(Save)
	if err := decoder.Decode(metadata); err != nil {
		log.Panicf("LoadWorld() - failed to decode metadata - %v", err)
	}

	log.Printf("LoadWorld() - loaded metadata; seed - %v", metadata.Seed)

	return NewWorld(metadata.Name, metadata.UUID, metadata.Seed)
}

// NOTE: world folder is named after the UUID, not after the world name
// that is, to avoid folder collision
func (world *World) Save() {
	saveDir := filepath.Join(config.WorldSaveDirectory, world.Metadata.UUID.String())

	// make a save directory, if it doesn't exist yet
	os.Mkdir(saveDir, os.ModePerm)

	// open world metadata file
	f, err := os.Create(filepath.Join(saveDir, config.WorldInfoFile))
	if err != nil {
		log.Panicf("failed to create world metadata file")
	}
	defer f.Close()

	// encode the metadata to it
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(world.Metadata); err != nil {
		log.Panicf("failed to encode world metadata")
	}

	// loop over all loaded chunks, saving modified ones to the disk
	for _, chunk := range world.chunks {
		chunk.Save(world.Metadata.UUID)
	}

	log.Println("World.Save() - saved")
}

func ChunkExistsOnDisk(id uuid.UUID, x, y uint64) bool {
	path := filepath.Join(config.WorldSaveDirectory, id.String(),
		fmt.Sprintf("chunk_%v_%v.gob", x, y))

	_, err := os.Stat(path)
	return err == nil
}

// if saved chunk doesn't exist, returns nil
func LoadChunk(id uuid.UUID, x, y uint64) *Chunk {
	path := filepath.Join(config.WorldSaveDirectory, id.String(),
		fmt.Sprintf("chunk_%v_%v.gob", x, y))

	if _, err := os.Stat(path); err == nil {
		f, err := os.Open(path)
		if err != nil {
			log.Panicf("failed to open a chunk - %v", err)
		}

		savedChunk := new(SavedChunk)
		if err := gob.NewDecoder(f).Decode(savedChunk); err != nil {
			log.Panicf("failed to decode a chunk - %v", err)
		}

		c := NewChunk(x, y)

		// decode blocks
		for x := uint(0); x < 16; x++ {
			for y := uint(0); y < 16; y++ {
				b := blocks.GetBlockByID(savedChunk.Data[x][y].Type)
				b.LoadState(savedChunk.Data[x][y].State)
				c.SetBlock(x, y, b)
			}
		}

		// mark chunk as unmodified, to avoid recursive loading/saving
		c.modified = false
		log.Printf("LoadChunk() - loaded chunk %v; %v from disk", x, y)
		return c
	}

	return nil
}

func (c *Chunk) Save(id uuid.UUID) {
	// if chunk wasn't modified, saving is unnecessary
	if !c.modified {
		return
	}

	path := filepath.Join(config.WorldSaveDirectory, id.String(),
		fmt.Sprintf("chunk_%v_%v.gob", c.x, c.y))

	f, err := os.Create(path)
	if err != nil {
		log.Panicf("failed to create chunk save file - %v", err)
	}
	defer f.Close()

	// serialize the chunk
	chunk := SavedChunk{
		X: c.x, Y: c.y,
	}
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			block := c.blocks[x][y]
			chunk.Data[x][y] = SavedBlock{
				Type:  block.Type(),
				State: block.State(),
			}
		}
	}

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(chunk); err != nil {
		log.Panicf("failed to encode chunk")
	}

	c.modified = false
	log.Printf("Chunk.Save() - %v; %v", c.x, c.y)
}

func DeleteWorld(id uuid.UUID) {
	path := filepath.Join(config.WorldSaveDirectory, id.String())
	if err := os.RemoveAll(path); err != nil {
		log.Panicf("Failed to delete world %v - %v", id, err)
	}
}

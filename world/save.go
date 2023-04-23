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
	Metadata types.Save

	saveRequests     chan Chunk
	loadRequestsPool map[types.Vec2u]bool
	// loadRequestsPool keeps track of currently requested chunks,
	// so that one same chunk can't be requested twice
	loadRequests chan types.Vec2u
	loaded       chan *Chunk
}

func NewWorldSaverLoader(metadata types.Save) *SaverLoader {
	return &SaverLoader{
		Metadata: metadata,

		saveRequests:     make(chan Chunk, 1024),
		loadRequestsPool: make(map[types.Vec2u]bool),
		loadRequests:     make(chan types.Vec2u, 256),
		loaded:           make(chan *Chunk),
	}
}

func (sl *SaverLoader) runSaver() {
	for {
		chunk := <-sl.saveRequests

		// FIXME: This is probably not very save
		chunk.Save(sl.Metadata)
	}
}

func (sl *SaverLoader) runLoader() {
	for {
		request := <-sl.loadRequests

		c := LoadChunk(sl.Metadata, request.X, request.Y)
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
	coords := types.Vec2u{X: cx, Y: cy}
	if sl.loadRequestsPool[coords] {
		return
	}
	sl.loadRequestsPool[coords] = true
	sl.loadRequests <- coords
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

func Load(baseID, id uuid.UUID) *World {
	saveDir := filepath.Join(config.WorldSaveDirectory, baseID.String(), id.String())

	f, err := os.Open(filepath.Join(saveDir, config.WorldInfoFile))
	if err != nil {
		log.Panicf("world.Load() - invalid file descriptor - %v", err)
	}

	decoder := gob.NewDecoder(f)
	metadata := new(types.Save)
	if err := decoder.Decode(metadata); err != nil {
		log.Panicf("world.Load() - failed to decode metadata - %v", err)
	}

	log.Printf("world.Load() - loaded metadata; seed - %v", metadata.Seed)

	return NewWorld(*metadata)
}

// NOTE: world folder is named after the UUID, not after the world name
// that is, to avoid folder collision
func (world *World) Save() {
	saveDir := filepath.Join(config.WorldSaveDirectory, world.metadata.BaseUUID.String(), world.metadata.UUID.String())

	// make a save directory, if it doesn't exist yet
	os.MkdirAll(saveDir, os.ModePerm)

	// open world metadata file
	f, err := os.Create(filepath.Join(saveDir, config.WorldInfoFile))
	if err != nil {
		log.Panicf("failed to create world metadata file")
	}
	defer f.Close()

	// encode the metadata to it
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(world.metadata); err != nil {
		log.Panicf("failed to encode world metadata")
	}

	// loop over all loaded chunks, saving modified ones to the disk
	for _, chunk := range world.chunks {
		chunk.Save(world.metadata)
	}
}

func ChunkExistsOnDisk(id uuid.UUID, x, y uint64) bool {
	path := filepath.Join(config.WorldSaveDirectory, id.String(),
		fmt.Sprintf("chunk_%v_%v.gob", x, y))

	_, err := os.Stat(path)
	return err == nil
}

// if saved chunk doesn't exist, returns nil
func LoadChunk(metadata types.Save, x, y uint64) *Chunk {
	path := filepath.Join(config.WorldSaveDirectory, metadata.BaseUUID.String(), metadata.UUID.String(),
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
		return c
	}

	return nil
}

func (c *Chunk) Save(metadata types.Save) {
	// if chunk wasn't modified, saving is unnecessary
	if !c.modified {
		return
	}

	path := filepath.Join(config.WorldSaveDirectory, metadata.BaseUUID.String(), metadata.UUID.String(),
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
}

func DeleteWorld(metadata types.Save) {
	path := filepath.Join(config.WorldSaveDirectory, metadata.BaseUUID.String())
	if err := os.RemoveAll(path); err != nil {
		log.Panicf("Failed to delete world %v - %v", metadata.BaseUUID, err)
	}
}

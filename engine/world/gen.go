/*
Functions related to world generation
*/
package world

import (
	"log"
	"math/rand"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/util"
	"github.com/aquilax/go-perlin"
)

// WorldGenerator maintains chunk generation queue
// Actual generation happens in separate goroutine,
// so we don't have any freezes on the main thread
type WorldGenerator struct {
	// Separate perlin noise generators for base blocks and vegetation/features
	baseGenerator      *perlin.Perlin
	secondaryGenerator *perlin.Perlin

	// requestsPool keeps track of currently requested chunks,
	// so that one same chunk can't be requested twice
	requestsPool map[util.Coords2i]bool
	requests     chan util.Coords2i
	generated    chan *Chunk
}

func NewWorldGenerator(seed int64) *WorldGenerator {
	// make a random generator using global world seed
	world := rand.New(rand.NewSource(seed))

	// and generate perlin noise seeds, using it
	var (
		baseSeed      = world.Int63()
		secondarySeed = world.Int63()
	)

	return &WorldGenerator{
		baseGenerator:      perlin.NewPerlin(2, 2, 16, baseSeed),
		secondaryGenerator: perlin.NewPerlin(2, 2, 16, secondarySeed),

		requestsPool: make(map[util.Coords2i]bool),
		// for some reason, without buffering, it hangs
		requests:  make(chan util.Coords2i, 128),
		generated: make(chan *Chunk, 128),
	}
}

func (g *WorldGenerator) run() {
	for {
		// listen for incoming requests
		req := <-g.requests

		c := NewChunk(req.X, req.Y)
		if err := c.Generate(g.baseGenerator, g.secondaryGenerator); err != nil {
			log.Panicf("WorldGenerator.run() - error while generating chunk - %v", err)
		}

		// FIXME: This is probably not very safe to pass pointers between goroutines
		g.generated <- c
	}
}

// starts chunk generator in separate goroutine
func (g *WorldGenerator) Run() {
	go g.run()
}

// Requests a chunk generation
// Chunk can be retrieved later through WorldGenerator.Receive()
func (g *WorldGenerator) Generate(cx, cy int64) {
	coords := util.Coords2i{X: cx, Y: cy}
	if g.requestsPool[coords] {
		return
	}
	g.requests <- coords
	g.requestsPool[coords] = true
}

// Returns newly generated chunk
// If none are pending, returns nil
func (g *WorldGenerator) Receive() *Chunk {
	select {
	case c := <-g.generated:
		// remove chunk from request pool
		delete(g.requestsPool, c.Coords())
		return c
	default:
		return nil
	}
}

// Features are regular random numbers, that can be used while generating blocks.
// Useful, when we need simple RNG, not perlin noise.
//
// They are reproducible, and unique for each block
type BlockFeatures struct {
	i1 int64
	u1 uint64
	f1 float64
}

// TODO: Optimize this
func makeFeatures(p *perlin.Perlin, bx, by int64) BlockFeatures {
	// We make new random generator, using block coordinates and perlin noise generator
	// This ensures that we get the same result every time, using same arguments
	seed := int64((height(p, float64(bx), float64(by), config.PerlinNoiseScaleFactor) / 2) * float64(1<<63))
	r := rand.New(rand.NewSource(seed))

	return BlockFeatures{
		i1: r.Int63(),
		u1: r.Uint64(),
		f1: r.Float64(),
	}
}

// returns values from 0 to 2
//
// x and y are world(block) coordinates
func height(gen *perlin.Perlin, x, y, scale float64) float64 {
	return gen.Noise2D(x/scale, y/scale) + 1
}

// generates basic blocks ( sand, water, etc. )
func genBase(baseGenerator *perlin.Perlin, x, y float64) Block {
	baseHeight := height(baseGenerator, x, y, config.PerlinNoiseScaleFactor)

	switch {
	case baseHeight <= 1: // Water
		return NewWaterBlock()
	case baseHeight <= 1.1: // Sand
		return NewSandBlock(false)
	case baseHeight <= 1.45: // Grass
		return NewGrassBlock()
	case baseHeight <= 1.65: // Stone
		return NewStoneBlock()
	default: // Snow
		return NewSnowBlock()
	}
}

// Checks if 8 neighbors of the block are of the same type
func checkNeighbors(desiredType BlockType, baseGenerator *perlin.Perlin, x, y float64) bool {
	sides := [8][2]float64{
		{x - 1, y},     // left
		{x + 1, y},     // right
		{x, y - 1},     // top
		{x, y + 1},     // bottom
		{x - 1, y - 1}, // top-left
		{x + 1, y - 1}, // top-right
		{x - 1, y + 1}, // bottom-left
		{x + 1, y + 1}, // bottom-right
	}

	for _, side := range sides {
		if genBase(baseGenerator, side[0], side[1]).Type() != desiredType {
			return false
		}
	}

	return true
}

// generates block features, depending on previous block
func genFeatures(previous Block, baseGenerator *perlin.Perlin, secondaryGenerator *perlin.Perlin, features BlockFeatures, x, y float64) Block {
	secondaryHeight := height(secondaryGenerator, x, y, config.PerlinNoiseScaleFactor)

	switch previous.Type() {
	case Sand:
		// With 3% change, generate sand with stones
		if features.f1 <= 0.03 {
			return NewSandBlock(true)
		}
	case Grass:
		// generate features on grass, only if it is surrounded by grass on all sides
		if !checkNeighbors(Grass, baseGenerator, x, y) {
			return previous
		}

		switch {
		case secondaryHeight <= 0.9: // Empty grass
			return previous
		case secondaryHeight <= 1.3: // Short grass / flowers
			// with 8% chance, make flowered grass
			if features.f1 <= 0.08 {
				return NewFlowersBlock()
			}
			return NewShortGrassBlock()
		default: // Tall grass
			return NewTallGrassBlock()
		}
	}

	// pass the base block forward, without any modifications
	return previous
}

// generates ground block at given coordinates
func gen(baseGenerator, secondaryGenerator *perlin.Perlin, x, y float64) Block {
	base := genBase(baseGenerator, x, y)
	withFeatures := genFeatures(base, baseGenerator, secondaryGenerator, makeFeatures(secondaryGenerator, int64(x), int64(y)), x, y)

	return withFeatures
}

func (c *Chunk) Generate(baseGenerator, secondaryGenerator *perlin.Perlin) error {
	// check if the chunk is out of the world borders
	// if it is, don't generate anything, instead return a chunk filled with stone blocks
	chunkOutOfBorders := c.x < 0 || c.y < 0 || c.x >= config.WorldWidth || c.y >= config.WorldHeight

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			bx := c.x*16 + int64(x)
			by := c.y*16 + int64(y)

			var generatedBlock Block = NewStoneBlock()
			if !chunkOutOfBorders {
				generatedBlock = gen(baseGenerator, secondaryGenerator, float64(bx), float64(by))
			}

			if err := c.SetBlock(x, y, generatedBlock); err != nil {
				return err
			}
		}
	}

	return nil
}

// simply fills a chunk with water
func (c *Chunk) GenerateDummy() error {
	// check if the chunk is out of the world borders
	// if it is, don't generate anything, instead return a chunk filled with stone blocks
	chunkOutOfBorders := c.x < 0 || c.y < 0 || c.x >= config.WorldWidth || c.y >= config.WorldHeight

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {

			var generatedBlock Block = NewStoneBlock()
			if !chunkOutOfBorders {
				generatedBlock = NewWaterBlock()
			}

			if err := c.SetBlock(x, y, generatedBlock); err != nil {
				return err
			}
		}
	}

	return nil
}

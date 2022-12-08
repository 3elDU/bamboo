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
	// Separate perlin noise generators for each layer
	bottomGenerator *perlin.Perlin
	groundGenerator *perlin.Perlin
	topGenerator    *perlin.Perlin

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
		bottomSeed = world.Int63()
		groundSeed = world.Int63()
		topSeed    = world.Int63()
	)

	return &WorldGenerator{
		bottomGenerator: perlin.NewPerlin(2, 2, 16, bottomSeed),
		groundGenerator: perlin.NewPerlin(2, 2, 16, groundSeed),
		topGenerator:    perlin.NewPerlin(2, 2, 16, topSeed),

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
		if err := c.Generate(g.bottomGenerator, g.groundGenerator, g.topGenerator); err != nil {
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
// Useful, when we need regular number generation, not perlin noise.
// For example, while generating vegetation, with 10% change we generate a flower instead of regular grass.
// So, we check block features, and depending on that, we decide.
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

// generates bottom block at given coordinates
func genBottom(p *perlin.Perlin, features BlockFeatures, x, y float64) Block {
	// h := height(source.p, x, y, config.PerlinNoiseScaleFactor/4)
	return NewEmptyBlock()
}

// generates ground block at given coordinates
func genGround(p *perlin.Perlin, prevBlock Block, features BlockFeatures, x, y float64) Block {
	// returns a value from 0 to 2
	h := height(p, x, y, config.PerlinNoiseScaleFactor)

	switch {
	case h <= 1: // Water
		return NewWaterBlock()
	case h <= 1.1: // Sand
		return NewSandBlock()
	case h <= 1.45: // Grass
		return NewGrassBlock()
	case h <= 1.65: // Stone
		return NewStoneBlock()
	default: // Snow
		return NewSnowBlock()
	}
}

// generates top block at given coordinates
func genTop(p *perlin.Perlin, prevBlock Block, features BlockFeatures, x, y float64) Block {
	h := height(p, x, y, config.PerlinNoiseScaleFactor/3)

	// check for the underlying block
	switch prevBlock.Type() {
	case Grass:
		var vegetationType BlockType

		switch {
		case h <= 0.9:
			return NewEmptyBlock()
		case h <= 1.3:
			// with 8% chance, make flowered grass
			if features.f1 <= 0.08 {
				vegetationType = Grass_Flowers
			} else {
				vegetationType = Grass_Plants_Small
			}
		default:
			vegetationType = Grass_Plants
		}

		return NewGrassVegetationBlock(vegetationType)
	case Sand:
		// With 3% change, generate stone with gravel
		if features.f1 <= 0.03 {
			return NewSandStoneBlock()
		} else {
			return NewEmptyBlock()
		}
	default:
		return NewEmptyBlock()
	}
}

func (c *Chunk) Generate(bottom, ground, top *perlin.Perlin) error {
	// log.Printf("Chunk.Generate() - Generating chunk at %v, %v", c.x, c.y)

	// check if the chunk is out of the world borders
	// if it is, don't generate a world, instead return a chunk filled with stone blocks
	chunkOutOfBorders := c.x < 0 || c.y < 0 || c.x >= config.WorldWidth || c.y >= config.WorldHeight

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			bx := c.x*16 + int64(x)
			by := c.y*16 + int64(y)

			var bottomBlock, groundBlock, topBlock Block

			if chunkOutOfBorders {
				bottomBlock = NewStoneBlock()
				groundBlock = NewEmptyBlock()
				topBlock = NewEmptyBlock()
			} else {
				bottomBlock = genBottom(bottom, makeFeatures(bottom, bx, by), float64(bx), float64(by))
				groundBlock = genGround(ground, bottomBlock, makeFeatures(bottom, bx, by), float64(bx), float64(by))
				topBlock = genTop(top, groundBlock, makeFeatures(bottom, bx, by), float64(bx), float64(by))
			}

			if err := c.SetStack(x, y, BlockStack{
				Bottom: bottomBlock,
				Ground: groundBlock,
				Top:    topBlock,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// simply fills a chunk with water
func (c *Chunk) GenerateDummy() error {
	// check if the chunk is out of the world borders
	// if it is, don't generate a world, instead return a chunk filled with stone blocks
	chunkOutOfBorders := c.x < 0 || c.y < 0 || c.x >= config.WorldWidth || c.y >= config.WorldHeight

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			var bottomBlock, groundBlock, topBlock Block

			if chunkOutOfBorders {
				bottomBlock = NewStoneBlock()
				groundBlock = NewEmptyBlock()
				topBlock = NewEmptyBlock()
			} else {
				bottomBlock = NewEmptyBlock()
				groundBlock = NewWaterBlock()
				topBlock = NewWaterBlock()
			}

			if err := c.SetStack(x, y, BlockStack{
				Bottom: bottomBlock,
				Ground: groundBlock,
				Top:    topBlock,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

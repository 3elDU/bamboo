/*
Functions related to world generation
*/
package world

import (
	"log"
	"math"
	"math/rand"

	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/aquilax/go-perlin"
)

// WorldGenerator maintains chunk generation queue
// Actual generation happens in separate goroutine,
// so we don't have any freezes on the main thread
type WorldGenerator struct {
	// Separate perlin noise generators for base blocks and vegetation/features
	baseGenerator      *perlin.Perlin
	secondaryGenerator *perlin.Perlin
	mountainGenerator  *perlin.Perlin

	// requestsPool keeps track of currently requested chunks,
	// so that one same chunk can't be requested twice
	requestsPool map[types.Coords2u]bool
	requests     chan types.Coords2u
	generated    chan *Chunk
}

func NewWorldGenerator(seed int64) *WorldGenerator {
	// make a random generator using global world seed
	world := rand.New(rand.NewSource(seed))

	// and generate perlin noise seeds, using it
	var (
		baseSeed      = world.Int63()
		secondarySeed = world.Int63()
		mountainSeed  = world.Int63()
	)

	return &WorldGenerator{
		baseGenerator:      perlin.NewPerlin(2, 2, 16, baseSeed),
		secondaryGenerator: perlin.NewPerlin(2, 2, 16, secondarySeed),
		mountainGenerator:  perlin.NewPerlin(2, 2, 16, mountainSeed),

		requestsPool: make(map[types.Coords2u]bool),
		// for some reason, without buffering, it hangs
		requests:  make(chan types.Coords2u, 128),
		generated: make(chan *Chunk, 128),
	}
}

func (g *WorldGenerator) run() {
	for {
		// listen for incoming requests
		req := <-g.requests

		c := NewChunk(req.X, req.Y)
		if err := c.Generate(g.baseGenerator, g.secondaryGenerator, g.mountainGenerator); err != nil {
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
func (g *WorldGenerator) Generate(cx, cy uint64) {
	coords := types.Coords2u{X: cx, Y: cy}
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
	f2 float64
}

// TODO: Optimize this
func makeFeatures(p *perlin.Perlin, bx, by uint64) BlockFeatures {
	// We make new random generator, using block coordinates and perlin noise generator
	// This ensures that we get the same result every time, using same arguments
	seed := int64((height(p, bx, by, config.PerlinNoiseScaleFactor) / 2) * float64(1<<63))
	r := rand.New(rand.NewSource(seed))

	return BlockFeatures{
		i1: r.Int63(),
		u1: r.Uint64(),
		f1: r.Float64(),
		f2: r.Float64(),
	}
}

// Applies circular mask to generated perlin noise
// The further block is from the center, the stronger the mask will be
// This makes the world look like an archipelago, surrounded by ocean on all sides,
// not like an infinite number of islands
func applyCircularMask(x, y float64, val float64) float64 {
	// TODO: make an "advanced options" sub-menu, in "new world" menu,
	// and move those constants there
	const (
		radius  = float64(config.WorldWidth) / 2.5
		centerX = float64(config.WorldWidth) / 2
		centerY = float64(config.WorldHeight) / 2
	)

	pointInsideCircle := math.Pow(float64(x)-centerX, 2)+math.Pow(float64(y)-centerY, 2) < math.Pow(radius, 2)
	if !pointInsideCircle {
		return 0
	}

	distanceToCenter := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
	// Divide the mask by 1.5, so it won't be too big
	mask := distanceToCenter / radius / 1.5
	return val - mask
}

// returns values from 0 to 2
//
// x and y are world(block) coordinates
func height(gen *perlin.Perlin, x, y uint64, scale float64) float64 {
	return gen.Noise2D(float64(x)/scale, float64(y)/scale) + 1
}

// generates basic blocks ( sand, water, etc. )
func genBase(baseGenerator *perlin.Perlin, x, y uint64) types.Block {
	baseHeight := applyCircularMask(float64(x), float64(y),
		height(baseGenerator, x, y, config.PerlinNoiseScaleFactor),
	)

	switch {
	case baseHeight <= 1: // Water
		return blocks.NewWaterBlock()
	case baseHeight <= 1.1: // Sand
		return blocks.NewSandBlock(false)
	default: // Grass
		return blocks.NewGrassBlock()
	}
}

func genMountain(previous types.Block, mountainGenerator *perlin.Perlin, x, y uint64) types.Block {
	if previous.Type() != blocks.Grass {
		return previous
	}

	mountainHeight := height(mountainGenerator, x/2, y/2, config.PerlinNoiseScaleFactor)

	if mountainHeight > 1.25 {
		return blocks.NewStoneBlock()
	}
	return previous
}

// Checks if 8 neighbors of the block are of the same type
func checkNeighbors(desiredType types.BlockType, baseGenerator *perlin.Perlin, x, y uint64) bool {
	sides := [8][2]uint64{
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
func genFeatures(previous types.Block, baseGenerator *perlin.Perlin, secondaryGenerator *perlin.Perlin, features BlockFeatures, x, y uint64) types.Block {
	// do not apply circular mask, while generating block features
	secondaryHeight := height(secondaryGenerator, x, y, config.PerlinNoiseScaleFactor)

	switch previous.Type() {
	case blocks.Sand:
		// With 3% change, generate sand with stones
		if features.f1 <= 0.03 {
			return blocks.NewSandBlock(true)
		}
	case blocks.Grass:
		// generate features on grass, only if it is surrounded by grass on all sides
		if !checkNeighbors(blocks.Grass, baseGenerator, x, y) {
			return previous
		}

		switch {
		case secondaryHeight <= 0.9: // Empty grass
			return previous
		case secondaryHeight <= 1.3: // Foliage
			switch {
			// with 1.5% chance, generate mushroom
			case features.f1 <= 0.015:
				if features.f2 <= 0.5 {
					return blocks.NewRedMushroomBlock()
				} else {
					return blocks.NewWhiteMushroomBlock()
				}
			// with 6% chance, generate flowers
			case features.f1 <= 0.06:
				return blocks.NewFlowersBlock()
			}

			return blocks.NewShortGrassBlock()
		default: // Tree
			return blocks.NewPineTreeBlock()
		}
	}

	// pass the base block forward, without any modifications
	return previous
}

// generates ground block at given coordinates
func gen(baseGenerator, secondaryGenerator, mountainGenerator *perlin.Perlin, x, y uint64) types.Block {
	var generated types.Block

	generated = genBase(baseGenerator, x, y)
	generated = genMountain(generated, mountainGenerator, x, y)
	generated = genFeatures(generated, baseGenerator, secondaryGenerator, makeFeatures(secondaryGenerator, x, y), x, y)

	return generated
}

func (c *Chunk) Generate(baseGenerator, secondaryGenerator, mountainGenerator *perlin.Perlin) error {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			bx := c.x*16 + uint64(x)
			by := c.y*16 + uint64(y)

			if err := c.SetBlock(x, y, gen(baseGenerator, secondaryGenerator, mountainGenerator, bx, by)); err != nil {
				return err
			}
		}
	}

	return nil
}

// simply fills a chunk with water
func (c *Chunk) GenerateDummy() error {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			if err := c.SetBlock(x, y, blocks.NewWaterBlock()); err != nil {
				return err
			}
		}
	}

	return nil
}

/*
Functions related to world generation
*/
package worldgen

import (
	"math"
	"math/rand"

	"github.com/3elDU/bamboo/blocks"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/aquilax/go-perlin"
)

// Chances of generating certain structures, blocks, and other worldgen-related constants
const (
	// Uses base height.
	// Height, below which water will generate
	WaterHeight = 1.0
	// Height, below which sand will generate
	SandHeight = 1.1

	// Uses secondary height.
	// Height, below which grass will generate
	GrassHeight = 0.9
	// Height, below which foliage will generate
	FoliageHeight = 1.3

	// %Chance of stones on sand being generated
	SandStoneChance = 0.03
	// %Chance of generating a mushroom
	MushroomChance = 0.015
	// %Chance of generating a flower
	FlowerChance = 0.06
	// %Chance of generating cave entrance in a chunk
	CaveEntranceChance = 0.05
)

// Generator maintains chunk generation queue.
// Actual generation happens in separate goroutine,
// so we don't have any freezes on the main thread
type Generator struct {
	// Separate perlin noise generators for base blocks and vegetation/features
	baseGenerator      *perlin.Perlin
	secondaryGenerator *perlin.Perlin

	// requestsPool keeps track of currently requested chunks,
	// so that one same chunk can't be requested twice
	requestsPool map[types.Vec2u]types.Chunk
	requests     chan types.Vec2u
	generated    chan types.Chunk
}

func NewWorldGenerator(seed int64) *Generator {
	// make a random generator using global world seed
	globalSeed := rand.New(rand.NewSource(seed))

	// generate perlin noise seeds, using it
	var (
		baseSeed      = globalSeed.Int63()
		secondarySeed = globalSeed.Int63()
	)

	return &Generator{
		baseGenerator:      perlin.NewPerlin(2, 2, 16, baseSeed),
		secondaryGenerator: perlin.NewPerlin(2, 2, 16, secondarySeed),

		requestsPool: make(map[types.Vec2u]types.Chunk),
		// for some reason, without buffering, it hangs
		requests:  make(chan types.Vec2u, 128),
		generated: make(chan types.Chunk, 128),
	}
}

func (generator *Generator) run() {
	for {
		// listen for incoming requests
		req := <-generator.requests

		chunk := generator.requestsPool[req]
		generator.generate(chunk)

		// FIXME: This is probably not very safe to pass pointers between goroutines
		generator.generated <- chunk
	}
}

// starts chunk generator in separate goroutine
func (generator *Generator) Run() {
	go generator.run()
}

// Requests a chunk generation
// Chunk can be retrieved later through Generator.Receive()
func (generator *Generator) Generate(chunk types.Chunk) {
	coords := chunk.Coords()
	if _, exists := generator.requestsPool[coords]; exists {
		return
	}
	generator.requests <- coords
	generator.requestsPool[coords] = chunk
}

// Returns newly generated chunk
// If none are pending, returns nil
func (generator *Generator) Receive() (chunks []types.Chunk) {
	noMoreValues := false
	for {
		select {
		case receivedChunk := <-generator.generated:
			chunks = append(chunks, receivedChunk)
			// remove chunk from request pool
			delete(generator.requestsPool, receivedChunk.Coords())
		default:
			noMoreValues = true
		}

		if noMoreValues {
			break
		}
	}
	return
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
	const (
		radius  = float64(config.WorldWidth) / 2.5
		centerX = float64(config.WorldWidth) / 2
		centerY = float64(config.WorldHeight) / 2
	)

	pointInsideCircle := math.Pow(x-centerX, 2)+math.Pow(y-centerY, 2) < math.Pow(radius, 2)
	if !pointInsideCircle {
		return 0
	}

	distanceToCenter := math.Sqrt(math.Pow(x-centerX, 2) + math.Pow(y-centerY, 2))
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
func (generator *Generator) genBase(x, y uint64) types.Block {
	baseHeight := applyCircularMask(float64(x), float64(y),
		height(generator.baseGenerator, x, y, config.PerlinNoiseScaleFactor),
	)

	switch {
	case baseHeight <= WaterHeight: // Water
		return blocks.NewWaterBlock()
	case baseHeight <= SandHeight: // Sand
		return blocks.NewSandBlock(false)
	default: // Grass
		return blocks.NewGrassBlock()
	}
}

// Checks if 8 neighbors of the block are of the same type
func (generator *Generator) checkNeighbors(desiredType types.BlockType, x, y uint64) bool {
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
		if generator.genBase(side[0], side[1]).Type() != desiredType {
			return false
		}
	}

	return true
}

// generates block features, depending on previous block
func (generator *Generator) genFeatures(previous types.Block, x, y uint64) types.Block {
	features := makeFeatures(generator.secondaryGenerator, x*16, y*16)

	// do not apply circular mask, while generating block features
	secondaryHeight := height(generator.secondaryGenerator, x, y, config.PerlinNoiseScaleFactor)

	switch previous.Type() {
	case blocks.Sand:
		// With 3% change, generate sand with stones
		if features.f1 <= SandStoneChance {
			return blocks.NewSandBlock(true)
		}
	case blocks.Grass:
		// generate features on grass, only if it is surrounded by grass on all sides
		if !generator.checkNeighbors(blocks.Grass, x, y) {
			return previous
		}

		switch {
		case secondaryHeight <= GrassHeight: // Empty grass
			return previous
		case secondaryHeight <= FoliageHeight: // Foliage
			switch {
			case features.f1 <= MushroomChance:
				if features.f2 <= 0.5 {
					return blocks.NewRedMushroomBlock()
				} else {
					return blocks.NewWhiteMushroomBlock()
				}
			case features.f1 <= FlowerChance:
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

func (generator *Generator) generateStructures(chunk types.Chunk) {
	chunkCoords := chunk.BlockCoords()
	features := makeFeatures(generator.secondaryGenerator, chunkCoords.X, chunkCoords.Y)

	if features.f1 < CaveEntranceChance {
		// use a bunch of hardcoded possible coordinates for cave entrance
		// it is just easier than trying to make reproducible RNG
		possibleCoordinates := []types.Vec2u{
			{X: 8, Y: 4},
			{X: 4, Y: 12},
			{X: 12, Y: 12},
		}

		// iterate over all possible coordinates, and pick the first valid pair
		var chosenCoordinates types.Vec2u
		valid := false
		for _, coords := range possibleCoordinates {
			// valid positions for cave entrance are those that are surrounded by grass blocks on all sides
			if generator.checkNeighbors(blocks.Grass, chunkCoords.X+coords.X, chunkCoords.Y+coords.Y) {
				chosenCoordinates = coords
				valid = true
				break
			}
		}
		if !valid {
			return
		}

		chunk.SetBlock(uint(chosenCoordinates.X), uint(chosenCoordinates.Y), blocks.NewCaveEntranceBlock())
	}
}

func (generator *Generator) generate(chunk types.Chunk) {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			chunkCoordinates := chunk.Coords()
			bx := chunkCoordinates.X*16 + uint64(x)
			by := chunkCoordinates.Y*16 + uint64(y)

			var generated types.Block
			generated = generator.genBase(bx, by)
			generated = generator.genFeatures(generated, bx, by)

			chunk.SetBlock(x, y, generated)
		}
	}

	generator.generateStructures(chunk)
}

// simply fills a chunk with water
func (generator *Generator) GenerateDummy(chunk types.Chunk) {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			chunk.SetBlock(x, y, blocks.NewWaterBlock())
		}
	}
}

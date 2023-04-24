package worldgen

import (
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/aquilax/go-perlin"
	"math"
	"math/rand"
)

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

type generatorImplementation interface {
	generate(chunk types.Chunk)
	generateDummy(chunk types.Chunk)
}

// Generator handles chunk generation queue, while the generation itself is handled by embedded class.
// Generation runs in separate goroutine to reduce freezes
type Generator struct {
	// requestsPool keeps track of currently requested chunks,
	// so that one same chunk can't be requested twice
	requestsPool map[types.Vec2u]types.Chunk
	requests     chan types.Vec2u
	generated    chan types.Chunk

	implementation generatorImplementation
}

func newGenerator(implementation generatorImplementation) *Generator {
	return &Generator{
		requestsPool: make(map[types.Vec2u]types.Chunk),
		// Use buffering for channels to be able to hold more than 1 request in a queue
		requests:  make(chan types.Vec2u, 1024),
		generated: make(chan types.Chunk, 1024),

		implementation: implementation,
	}
}

func (generator *Generator) Run() {
	for {
		request := <-generator.requests

		chunk := generator.requestsPool[request]
		generator.implementation.generate(chunk)

		generator.generated <- chunk
	}
}

// Requests a chunk generation
// Generated chunks can be received through generator.Receive()
func (generator *Generator) Generate(chunk types.Chunk) {
	coords := chunk.Coords()
	if _, exists := generator.requestsPool[coords]; exists {
		return
	}
	generator.requests <- coords
	generator.requestsPool[coords] = chunk
}

// Unlike Generator.Generate(), Generator.GenerateDummy() generates the chunk immediately, without the queue
func (generator *Generator) GenerateDummy(chunk types.Chunk) {
	generator.implementation.generateDummy(chunk)
}

// Returns a list of generated chunks
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

/*
Functions related to world generation
*/
package world

import (
	"math/rand"

	"github.com/3elDU/bamboo/config"
	"github.com/aquilax/go-perlin"
)

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
	return NewStoneBlock(0)
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
		return NewStoneBlock(h)
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
				bottomBlock = NewStoneBlock(0)
				groundBlock = NewEmptyBlock()
				topBlock = NewEmptyBlock()
			} else {
				bottomBlock = genBottom(bottom, makeFeatures(bottom, bx, by), float64(bx), float64(by))
				groundBlock = genGround(ground, bottomBlock, makeFeatures(bottom, bx, by), float64(bx), float64(by))
				topBlock = genTop(top, groundBlock, makeFeatures(bottom, bx, by), float64(bx), float64(by))
			}

			err := c.SetStack(x, y, BlockStack{
				bottom: bottomBlock,
				ground: groundBlock,
				top:    topBlock,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

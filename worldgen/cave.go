package worldgen

import (
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/aquilax/go-perlin"
)

type CaveGenerator struct {
	noiseSeed int64
	noise     *perlin.Perlin
}

func NewCaveGenerator(seed int64) types.WorldGenerator {
	implementation := &CaveGenerator{
		noiseSeed: seed,
		noise:     perlin.NewPerlin(2, 2, 1, seed),
	}
	return newGenerator(implementation)
}

func (generator *CaveGenerator) generate(chunk types.Chunk) {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			coords := chunk.BlockCoords()
			h := height(generator.noise, coords.X+uint64(x), coords.Y+uint64(y), config.PerlinNoiseScaleFactor/5)

			if h < 1 {
				chunk.SetBlock(x, y, types.NewCaveFloorBlock())
			} else {
				chunk.SetBlock(x, y, types.NewCaveWallBlock())
			}
		}
	}
}

func (generator *CaveGenerator) generateDummy(chunk types.Chunk) {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			chunk.SetBlock(x, y, types.NewCaveFloorBlock())
		}
	}
}

func (generator *CaveGenerator) seed() int64 {
	return generator.noiseSeed
}

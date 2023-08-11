package worldgen

import (
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/aquilax/go-perlin"
)

const (
	// Chance that cave floor will have little sprouts growing on it
	CaveFloorSproutsChance = 0.1

	// %Chance that iron ore will be generated on this block
	// (Iron ore can only generated on cave floor block, not inside walls)
	IronOreChance = 0.006
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

func (generator *CaveGenerator) generateBase(x, y uint64) types.Block {
	h := height(generator.noise, x, y, config.PerlinNoiseScaleFactor/16)
	features := makeFeatures(generator.noise, x, y)

	if h < 1 {
		return types.NewCaveFloorBlock(
			features.f1 < CaveFloorSproutsChance,
		)
	} else {
		return types.NewCaveWallBlock()
	}
}

// checks if block is neighboring with air (cave floor) block on all four sides
func (generator *CaveGenerator) blockSurroundedByAir(x, y uint64) bool {
	neighbors := []types.Vec2u{
		{X: x, Y: y - 1}, // top
		{X: x - 1, Y: y}, // left
		{X: x + 1, Y: y}, // right
		{X: x, Y: y + 1}, // bottom
	}

	for _, neighbor := range neighbors {
		if generator.generateBase(neighbor.X, neighbor.Y).Type() != types.CaveFloorBlock {
			return false
		}
	}
	return true
}

func (generator *CaveGenerator) generateOre(x, y uint64, previous types.Block) types.Block {
	if previous.Type() != types.CaveFloorBlock {
		return previous
	}

	features := makeFeatures(generator.noise, x, y)

	// if a block is surrounded by air, chances of generating ore blocks are 2 times lower
	surroundedByAir := generator.blockSurroundedByAir(x, y)
	var chanceDivider = 1.0
	if surroundedByAir {
		chanceDivider = 2.0
	}

	if features.f1 < IronOreChance/chanceDivider {
		return types.NewIronOreBlock()
	}
	return previous
}

func (generator *CaveGenerator) generate(chunk types.Chunk) {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			coords := chunk.BlockCoords()
			// Block coordinates in world space
			bx, by := coords.X+uint64(x), coords.Y+uint64(y)

			block := generator.generateBase(bx, by)
			block = generator.generateOre(bx, by, block)
			chunk.SetBlock(x, y, block)
		}
	}
}

func (generator *CaveGenerator) generateDummy(chunk types.Chunk) {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			chunk.SetBlock(x, y, types.NewCaveFloorBlock(false))
		}
	}
}

func (generator *CaveGenerator) seed() int64 {
	return generator.noiseSeed
}

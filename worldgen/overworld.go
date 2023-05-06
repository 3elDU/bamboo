/*
Functions related to world generation
*/
package worldgen

import (
	"github.com/google/uuid"
	"log"
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

type OverworldGenerator struct {
	metadata types.Save
	// Separate perlin noise generators for base blocks and vegetation/features
	basePerlin      *perlin.Perlin
	secondaryPerlin *perlin.Perlin
}

func NewOverworldGenerator(metadata types.Save) types.WorldGenerator {
	// make a random generator using global world seed
	globalSeed := rand.New(rand.NewSource(metadata.Seed))

	// generate perlin noise seeds, using it
	var (
		baseSeed      = globalSeed.Int63()
		secondarySeed = globalSeed.Int63()
	)

	implementation := &OverworldGenerator{
		metadata:        metadata,
		basePerlin:      perlin.NewPerlin(2, 2, 16, baseSeed),
		secondaryPerlin: perlin.NewPerlin(2, 2, 16, secondarySeed),
	}

	return newGenerator(implementation)
}

// generates basic blocks ( sand, water, etc. )
func (generator *OverworldGenerator) genBase(x, y uint64) types.Block {
	baseHeight := applyCircularMask(generator.metadata.Size, float64(x), float64(y),
		height(generator.basePerlin, x, y, config.PerlinNoiseScaleFactor),
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
func (generator *OverworldGenerator) checkNeighbors(desiredType types.BlockType, x, y uint64) bool {
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
func (generator *OverworldGenerator) genFeatures(previous types.Block, x, y uint64) types.Block {
	features := makeFeatures(generator.secondaryPerlin, x*16, y*16)

	// do not apply circular mask, while generating block features
	secondaryHeight := height(generator.secondaryPerlin, x, y, config.PerlinNoiseScaleFactor)

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

func (generator *OverworldGenerator) generateStructures(chunk types.Chunk) {
	chunkCoords := chunk.BlockCoords()
	features := makeFeatures(generator.secondaryPerlin, chunkCoords.X, chunkCoords.Y)

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

		// kinda slow but reproducible with the same seed, which is the most important
		rng := rand.New(rand.NewSource(features.i1))
		id, err := uuid.NewRandomFromReader(rng)
		if err != nil {
			// this should really never happen
			log.Panicf("failed to generate UUID for cave: %v", err)
		}
		log.Printf("cave at %v, %v: %v", chunkCoords.X+chosenCoordinates.X, chunkCoords.Y+chosenCoordinates.Y, id)
		chunk.SetBlock(uint(chosenCoordinates.X), uint(chosenCoordinates.Y), blocks.NewCaveEntranceBlock(id))
	}
}

func (generator *OverworldGenerator) generate(chunk types.Chunk) {
	chunkCoordinates := chunk.Coords()

	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
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
func (generator *OverworldGenerator) generateDummy(chunk types.Chunk) {
	for x := uint(0); x < 16; x++ {
		for y := uint(0); y < 16; y++ {
			chunk.SetBlock(x, y, blocks.NewWaterBlock())
		}
	}
}

func (generator *OverworldGenerator) seed() int64 {
	return generator.metadata.Seed
}

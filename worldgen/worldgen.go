package worldgen

import (
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/world_type"
	"log"
)

// Returns a world generator for that specific world type
func NewWorldgenForWorld(metadata types.Save) types.WorldGenerator {
	switch metadata.WorldType {
	case world_type.Overworld:
		return NewOverworldGenerator(metadata)
	case world_type.Cave:
		return NewCaveGenerator(metadata.Seed)
	}

	log.Panicf("no world generator associated with world type %v", metadata)
	return nil
}

package worldgen

import (
	"github.com/3elDU/bamboo/types"
	"github.com/3elDU/bamboo/world_type"
	"log"
)

// Returns a world generator for that specific world type
func NewWorldgenForType(seed int64, worldType world_type.WorldType) types.WorldGenerator {
	switch worldType {
	case world_type.Overworld:
		return NewOverworldGenerator(seed)
	case world_type.Cave:
		return NewCaveGenerator(seed)
	}

	log.Panicf("no world generator associated with world type %v", worldType)
	return nil
}

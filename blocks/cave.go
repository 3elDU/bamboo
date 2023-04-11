package blocks

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
	"github.com/google/uuid"
	"log"
)

func init() {
	gob.Register(CaveEntranceState{})
}

type CaveEntranceState struct {
	BaseBlockState
	TexturedBlockState
	ID uuid.UUID
}

type CaveEntranceBlock struct {
	baseBlock
	texturedBlock

	id uuid.UUID
}

func NewCaveEntranceBlock() *CaveEntranceBlock {
	return &CaveEntranceBlock{
		baseBlock: baseBlock{
			blockType: CaveEntrance,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("cave"),
		},
		id: uuid.New(),
	}
}

func (cave *CaveEntranceBlock) State() interface{} {
	return CaveEntranceState{
		BaseBlockState:     cave.baseBlock.State().(BaseBlockState),
		TexturedBlockState: cave.texturedBlock.State().(TexturedBlockState),
		ID:                 cave.id,
	}
}

func (cave *CaveEntranceBlock) LoadState(s interface{}) {
	state := s.(CaveEntranceState)
	cave.baseBlock.LoadState(state.BaseBlockState)
	cave.texturedBlock.LoadState(state.TexturedBlockState)
	cave.id = state.ID
}

func (cave *CaveEntranceBlock) Interact(world types.World, _ types.Coords2f) {
	log.Println("interacted!")
	world.ChunkAtB(uint64(cave.x), uint64(cave.y)).SetBlock(cave.x%16, cave.y%16, NewGrassBlock())
}

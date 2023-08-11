package blocks_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/event"
	"github.com/3elDU/bamboo/types"
	"github.com/google/uuid"
)

func init() {
	gob.Register(CaveEntranceState{})
	types.NewCaveEntranceBlock = NewCaveEntranceBlock
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

func NewCaveEntranceBlock(id uuid.UUID) types.Block {
	return &CaveEntranceBlock{
		baseBlock: baseBlock{
			blockType: types.CaveEntranceBlock,
		},
		texturedBlock: texturedBlock{
			tex: assets.Texture("cave"),
		},
		id: id,
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

func (cave *CaveEntranceBlock) Collide(_ types.World, _ types.Vec2f) {
	event.FireEvent(event.NewEvent(
		event.CaveEnter,
		event.CaveEnteredArgs{
			ID: cave.id,
		},
	))
}

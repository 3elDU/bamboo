package blocks_impl

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/event"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(CaveExitState{})
	types.NewCaveExitBlock = NewCaveExitBlock
}

type CaveExitState struct {
	BaseBlockState
	TexturedBlockState
}

type CaveExitBlock struct {
	baseBlock
	texturedBlock
}

func NewCaveExitBlock() types.Block {
	return &CaveExitBlock{
		baseBlock: baseBlock{
			blockType: types.CaveExitBlock,
		},
		texturedBlock: texturedBlock{
			tex: asset_loader.Texture("cave_exit"),
		},
	}
}

func (cave *CaveExitBlock) State() interface{} {
	return CaveEntranceState{
		BaseBlockState:     cave.baseBlock.State().(BaseBlockState),
		TexturedBlockState: cave.texturedBlock.State().(TexturedBlockState),
	}
}

func (cave *CaveExitBlock) LoadState(s interface{}) {
	state := s.(CaveEntranceState)
	cave.baseBlock.LoadState(state.BaseBlockState)
	cave.texturedBlock.LoadState(state.TexturedBlockState)
}

func (cave *CaveExitBlock) Interact(_ types.World, _ types.Vec2f) {
	event.FireEvent(event.NewEvent(
		event.CaveExit, nil,
	))
}

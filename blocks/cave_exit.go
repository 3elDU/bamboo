package blocks

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/event"
	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(CaveExitState{})
}

type CaveExitState struct {
	BaseBlockState
	TexturedBlockState
}

type CaveExitBlock struct {
	baseBlock
	texturedBlock
}

func NewCaveExitBlock() *CaveExitBlock {
	return &CaveExitBlock{
		baseBlock: baseBlock{
			blockType: CaveExit,
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

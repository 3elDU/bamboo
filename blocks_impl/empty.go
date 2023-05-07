/*
	Empty block
*/

package blocks_impl

import (
	"github.com/3elDU/bamboo/types"
)

func init() {
	types.NewEmptyBlock = NewEmptyBlock
}

type EmptyBlock struct {
	baseBlock
}

func (e *EmptyBlock) Update(_ types.World) {

}

func NewEmptyBlock() types.Block {
	return &EmptyBlock{
		baseBlock: baseBlock{
			blockType: types.EmptyBlock,
		},
	}
}

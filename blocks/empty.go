/*
	Empty block
*/

package blocks

import (
	"github.com/3elDU/bamboo/types"
)

type EmptyBlock struct {
	baseBlock
}

func (e *EmptyBlock) Update(_ types.World) {

}

func NewEmptyBlock() *EmptyBlock {
	return &EmptyBlock{
		baseBlock: baseBlock{
			blockType: Empty,
		},
	}
}

/*
	Empty block
*/

package blocks

import (
	"github.com/3elDU/bamboo/types"
)

type emptyBlock struct {
	baseBlock
}

func (e *emptyBlock) Update(_ types.World) {

}

func NewEmptyBlock() *emptyBlock {
	return &emptyBlock{
		baseBlock: baseBlock{
			blockType: Empty,
		},
	}
}

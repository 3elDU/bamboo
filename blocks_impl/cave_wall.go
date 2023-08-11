package blocks_impl

import (
	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/assets"
)

func init() {
	types.NewCaveWallBlock = NewCaveWallBlock
}

type CaveWallBlock struct {
	connectedBlock
	collidableBlock
}

func NewCaveWallBlock() types.Block {
	return &CaveWallBlock{
		connectedBlock: connectedBlock{
			baseBlock: baseBlock{
				blockType: types.CaveWallBlock,
			},
			tex:        assets.ConnectedTexture("cave_wall", false, false, false, false),
			connectsTo: []types.BlockType{types.CaveWallBlock},
		},
		collidableBlock: collidableBlock{
			collidable:      true,
			collisionPoints: defaultCollisionPoints(),
		},
	}
}

func (b *CaveWallBlock) State() interface{} {
	return nil
}

func (b *CaveWallBlock) LoadState(s interface{}) {

}

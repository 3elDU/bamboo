package blocks_impl

import (
	"encoding/gob"
	"log"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"
)

func init() {
	gob.Register(ConnectedBlockState{})
}

type ConnectedBlockState struct {
	BaseBlockState
	Texture        string
	SidesConnected [4]bool
}

type connectedBlock struct {
	baseBlock
	connectsTo     []types.BlockType
	tex            types.ConnectedTexture
	sidesConnected [4]bool
}

func (b *connectedBlock) shouldConnect(other types.BlockType) bool {
	return slices.Contains(b.connectsTo, other)
}

func (b *connectedBlock) Render(world types.World, screen *ebiten.Image, pos types.Vec2f) {
	var connectedSides [4]bool
	for i, side := range [4]types.Vec2i{{X: -1, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1}} {
		x, y := int64(b.x)+side.X, int64(b.y)+side.Y
		neighbor := world.BlockAt(uint64(x), uint64(y))
		if !b.shouldConnect(neighbor.Type()) {
			continue
		}

		connectedSides[i] = true
		// If neighbor is on another chunk, trigger redraw of that chunk
		if neighbor.ParentChunk() != b.parentChunk {
			neighbor.ParentChunk().TriggerRedraw()
		}
	}
	b.tex.SetConnectedSides(connectedSides)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(b.tex.Texture(), opts)
}

func (b *connectedBlock) TextureName() string {
	return b.tex.Name()
}

func (b *connectedBlock) State() interface{} {
	return ConnectedBlockState{
		BaseBlockState: b.baseBlock.State().(BaseBlockState),
		Texture:        b.tex.Name(),
		SidesConnected: b.sidesConnected,
	}
}

func (b *connectedBlock) LoadState(s interface{}) {
	state, ok := s.(ConnectedBlockState)
	if !ok {
		log.Panicf("%T - invalid state type; expected %T, got %T", b, ConnectedBlockState{}, state)
	}

	b.baseBlock.LoadState(state.BaseBlockState)
	b.tex = asset_loader.ConnectedTextureFromArray(state.Texture, state.SidesConnected)
}

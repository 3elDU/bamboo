package blocks

import (
	"encoding/gob"
	"log"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/3elDU/bamboo/texture"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"
)

func init() {
	gob.Register(ConnectedTextureState{})
	gob.Register(ConnectedBlockState{})
}

type ConnectedTextureState struct {
	Base  string
	Sides [4]bool
}

type ConnectedBlockState struct {
	BaseBlockState
	ConnectedTextureState
}

type connectedBlock struct {
	baseBlock
	connectsTo []types.BlockType
	tex        texture.ConnectedTexture
}

func (b *connectedBlock) shouldConnect(other types.BlockType) bool {
	return slices.Contains(b.connectsTo, other)
}

func (b *connectedBlock) Render(world types.World, screen *ebiten.Image, pos types.Coords2f) {
	var sidesConnected [4]bool
	for i, side := range [4]types.Coords2i{{X: -1, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1}} {
		x, y := int64(b.x)+side.X, int64(b.y)+side.Y
		neighbor := world.BlockAt(uint64(x), uint64(y))
		if !b.shouldConnect(neighbor.Type()) {
			continue
		}

		sidesConnected[i] = true
		// If neighbor is on another chunk, trigger redraw of that chunk
		if neighbor.ParentChunk() != b.parentChunk {
			neighbor.ParentChunk().TriggerRedraw()
		}
	}
	b.tex.SetSidesArray(sidesConnected)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(asset_loader.Texture(b.tex.FullName()).Texture(), opts)
}

func (b *connectedBlock) TextureName() string {
	return b.tex.FullName()
}

func (b *connectedBlock) State() interface{} {
	return ConnectedBlockState{
		BaseBlockState:        b.baseBlock.State().(BaseBlockState),
		ConnectedTextureState: ConnectedTextureState{Base: b.tex.Base, Sides: b.tex.SidesConnected},
	}
}

func (b *connectedBlock) LoadState(s interface{}) {
	state, ok := s.(ConnectedBlockState)
	if !ok {
		log.Panicf("%T - invalid state type; expected %T, got %T", b, ConnectedBlockState{}, state)
	}

	b.baseBlock.LoadState(state.BaseBlockState)
	b.tex = asset_loader.ConnectedTexture(state.ConnectedTextureState.Base,
		state.ConnectedTextureState.Sides[0],
		state.ConnectedTextureState.Sides[1],
		state.ConnectedTextureState.Sides[2],
		state.ConnectedTextureState.Sides[3],
	)
}

/*
	Declarations of basic block types
*/

package world

import (
	"encoding/gob"
	"fmt"
	"math"

	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/texture"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"
)

func init() {
	// Register types for proper serialization
	gob.Register(BaseBlockState{})
	gob.Register(TexturedBlockState{})
	gob.Register(CompositeBlockState{})
	gob.Register(ConnectedBlockState{})
	gob.Register(CollidableBlockState{})
}

type BlockType int

type Block interface {
	Coords() util.Coords2u
	SetCoords(coords util.Coords2u)
	ParentChunk() *Chunk
	SetParentChunk(chunk *Chunk)
	Type() BlockType

	State() interface{}
	LoadState(interface{}) error
}

type CollidableBlock interface {
	Block
	CollisionPoints() [4]util.Coords2f
	PlayerSpeed() float64
}

type DrawableBlock interface {
	Block
	Render(world *World, screen *ebiten.Image, pos util.Coords2f)
	TextureName() string
}

type UpdateableBlock interface {
	Block
	Update(world *World)
}

type BaseBlockState struct {
	BlockType BlockType
}

// Base structure inherited by all blocks
// Contains some basic parameters, so we don't have to implement them for ourselves
type baseBlock struct {
	// Usually you don't have to set this for youself,
	// Since world.Gen() sets them automatically
	parentChunk *Chunk
	// Block coordinates in world space
	x, y uint

	// Block types are defined in (blocks.go):13
	// Each block must specify it's type, so that we can actually know what the block it is
	// ( Remember, all blocks are the same interface )
	blockType BlockType
}

type TexturedBlockState struct {
	Name     string
	Rotation float64
}

// Another base structure, to simplify things
type texturedBlock struct {
	tex      texture.Texture
	rotation float64 // in degrees
}

type CompositeBlockState struct {
	BaseBlockState
	TexturedBlockState
}

type compositeBlock struct {
	baseBlock
	texturedBlock
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
	connectsTo []BlockType
	tex        texture.ConnectedTexture
}

type CollidableBlockState struct {
	CollisionPoints [4]util.Coords2f
	PlayerSpeed     float64
}

type collidableBlock struct {
	// Each collision point is local coordinate, where (0, 0) is block's top-left corner
	collisionPoints [4]util.Coords2f

	// How fast player could move through this block
	// Calculated by basePlayerSpeed * playerSpeed
	// Applicable only if collidable is false
	playerSpeed float64
}

func (b *baseBlock) Coords() util.Coords2u {
	return util.Coords2u{X: uint64(b.x), Y: uint64(b.y)}
}

func (b *baseBlock) SetCoords(coords util.Coords2u) {
	b.x = uint(coords.X)
	b.y = uint(coords.Y)
}

func (b *baseBlock) ParentChunk() *Chunk {
	return b.parentChunk
}

func (b *baseBlock) SetParentChunk(c *Chunk) {
	b.parentChunk = c
}

func (b *baseBlock) Type() BlockType {
	return b.blockType
}

func (b *baseBlock) State() interface{} {
	return BaseBlockState{
		BlockType: b.blockType,
	}
}

func (b *baseBlock) LoadState(state interface{}) error {
	if state, ok := state.(BaseBlockState); ok {
		b.blockType = state.BlockType
	} else {
		return fmt.Errorf("%T - invalid state type; expected %T, got %T", b, BaseBlockState{}, state)
	}
	return nil
}

func (b *texturedBlock) Render(_ *World, screen *ebiten.Image, pos util.Coords2f) {
	opts := &ebiten.DrawImageOptions{}

	if b.rotation != 0 {
		w, h := b.tex.Texture().Size()
		// Move image half a texture size, so that rotation origin will be in the center
		opts.GeoM.Translate(float64(-w/2), float64(-h/2))
		opts.GeoM.Rotate(b.rotation * (math.Pi / 180))
		pos.X += float64(w / 2)
		pos.Y += float64(h / 2)
	}

	opts.GeoM.Translate(pos.X, pos.Y)

	screen.DrawImage(b.tex.Texture(), opts)
}

func (b *texturedBlock) TextureName() string {
	return b.tex.Name
}

func (b *texturedBlock) State() interface{} {
	return TexturedBlockState{
		Name:     b.tex.Name,
		Rotation: b.rotation,
	}
}

func (b *texturedBlock) LoadState(state interface{}) error {
	if state, ok := state.(TexturedBlockState); ok {
		b.tex = asset_loader.Texture(state.Name)
		b.rotation = state.Rotation
	} else {
		return fmt.Errorf("%T - invalid state type; expected %T, got %T", b, TexturedBlockState{}, state)
	}
	return nil
}

// Dummy update method
func (b *compositeBlock) Update(world *World) {

}

func (b *compositeBlock) State() interface{} {
	return CompositeBlockState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		TexturedBlockState: b.texturedBlock.State().(TexturedBlockState),
	}
}

func (b *compositeBlock) LoadState(state interface{}) error {
	if state, ok := state.(CompositeBlockState); ok {
		if err := b.baseBlock.LoadState(state.BaseBlockState); err != nil {
			return err
		}
		if err := b.texturedBlock.LoadState(state.TexturedBlockState); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%T - invalid state type; expected %T, got %T", b, CompositeBlockState{}, state)
	}
	return nil
}

// Update connected sides
func (b *connectedBlock) Update(world *World) {

}

func (b *connectedBlock) shouldConnect(other BlockType) bool {
	return slices.Contains(b.connectsTo, other)
}

func (b *connectedBlock) Render(world *World, screen *ebiten.Image, pos util.Coords2f) {
	var sidesConnected [4]bool
	for i, side := range [4]util.Coords2i{{X: -1, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1}} {
		x, y := int64(b.x)+side.X, int64(b.y)+side.Y

		neighbor := world.BlockAt(uint64(x), uint64(y))
		if !b.shouldConnect(neighbor.Type()) {
			continue
		}

		sidesConnected[i] = true
		// If neighbor is on another chunk, trigger redraw of that chunk
		if neighbor.ParentChunk() != b.parentChunk {
			neighbor.ParentChunk().needsRedraw = true
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

func (b *connectedBlock) LoadState(state interface{}) error {
	if state, ok := state.(ConnectedBlockState); ok {
		if err := b.baseBlock.LoadState(state.BaseBlockState); err != nil {
			return err
		}
		b.tex = asset_loader.ConnectedTexture(state.ConnectedTextureState.Base,
			state.ConnectedTextureState.Sides[0],
			state.ConnectedTextureState.Sides[1],
			state.ConnectedTextureState.Sides[2],
			state.ConnectedTextureState.Sides[3],
		)
	} else {
		return fmt.Errorf("%T - invalid state type; expected %T, got %T", b, CompositeBlockState{}, state)
	}
	return nil
}

func (b *collidableBlock) PlayerSpeed() float64 {
	return b.playerSpeed
}

func (b *collidableBlock) CollisionPoints() [4]util.Coords2f {
	return b.collisionPoints
}

/*
	Declarations of basic block types
*/

package world

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"

	"github.com/3elDU/bamboo/engine/asset_loader"
	"github.com/3elDU/bamboo/engine/texture"
	"github.com/3elDU/bamboo/util"
	"github.com/hajimehoshi/ebiten/v2"
)

type Layer int

const (
	BottomLayer Layer = iota
	GroundLayer
	TopLayer
)

type BlockType int

type Block interface {
	Coords() util.Coords2i
	SetCoords(coords util.Coords2i)
	ParentChunk() *Chunk
	SetParentChunk(chunk *Chunk)
	SetLayer(layer Layer)
	Layer() Layer
	Type() BlockType

	// Whether the player should collide with the block
	Collidable() bool

	Update(world *World)
	Render(world *World, screen *ebiten.Image, pos util.Coords2f)
	TextureName() string

	State() []byte
	LoadState(state_bytes []byte) error
}

// The map is actually 3D, and consists of three layers:
//
//  1. Fossils / Ore
//  2. Ground block ( the one you`ll see the most )
//  3. Top block - decoration / vegetation / player buildings / etc.
type BlockStack struct {
	Bottom Block
	Ground Block
	Top    Block
}

type BaseBlockState struct {
	Collidable  bool
	PlayerSpeed float64
	BlockType   BlockType
}

// Base structure inherited by all blocks
// Contains some basic parameters, so we don't have to implement them for ourselves
type baseBlock struct {
	// Usually you don't have to set this for youself,
	// Since world.Gen() sets them automatically
	parentChunk *Chunk
	// Block coordinates in world space
	x, y  int
	layer Layer

	// Whether collision will work with this block
	collidable bool

	// How fast player could move through this block
	// Calculated by basePlayerSpeed * playerSpeed
	// Applicable only if collidable is false
	playerSpeed float64

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
	BaseBlockState     []byte
	TexturedBlockState []byte
}

type compositeBlock struct {
	baseBlock
	texturedBlock
}

type ConnectedTextureState struct {
	Base           string
	SidesConnected [4]bool
}

type ConnectedBlockState struct {
	BaseBlockState        []byte
	ConnectedTextureState []byte
}

type connectedBlock struct {
	baseBlock
	tex texture.ConnectedTexture
}

func (b *baseBlock) Coords() util.Coords2i {
	return util.Coords2i{X: int64(b.x), Y: int64(b.y)}
}

func (b *baseBlock) SetCoords(coords util.Coords2i) {
	b.x = int(coords.X)
	b.y = int(coords.Y)
}

func (b *baseBlock) ParentChunk() *Chunk {
	return b.parentChunk
}

func (b *baseBlock) SetParentChunk(c *Chunk) {
	b.parentChunk = c
}

func (b *baseBlock) Layer() Layer {
	return b.layer
}

func (b *baseBlock) SetLayer(layer Layer) {
	b.layer = layer
}

func (b *baseBlock) Collidable() bool {
	return b.collidable
}

func (b *baseBlock) Type() BlockType {
	return b.blockType
}

func (b *baseBlock) State() []byte {
	buf := new(bytes.Buffer)
	state := BaseBlockState{
		Collidable:  b.collidable,
		PlayerSpeed: b.playerSpeed,
	}
	if err := gob.NewEncoder(buf).Encode(state); err != nil {
		log.Panicf("baseBlock.State() - %v", err)
	}
	return buf.Bytes()
}

func (b *baseBlock) LoadState(state_bytes []byte) error {
	buf := bytes.NewBuffer(state_bytes)
	state := new(BaseBlockState)
	if err := gob.NewDecoder(buf).Decode(state); err != nil {
		return err
	}
	b.collidable = state.Collidable
	b.playerSpeed = state.PlayerSpeed
	return nil
}

func (b *texturedBlock) Render(_ *World, screen *ebiten.Image, pos util.Coords2f) {
	opts := &ebiten.DrawImageOptions{}

	if b.rotation != 0 {
		w, h := b.tex.Texture.Size()
		// Move image half a texture size, so that rotation origin will be in the center
		opts.GeoM.Translate(float64(-w/2), float64(-h/2))
		opts.GeoM.Rotate(b.rotation * (math.Pi / 180))
		pos.X += float64(w / 2)
		pos.Y += float64(h / 2)
	}

	opts.GeoM.Translate(pos.X, pos.Y)

	screen.DrawImage(b.tex.Texture, opts)
}

func (b *texturedBlock) TextureName() string {
	return b.tex.Name
}

func (b *texturedBlock) State() []byte {
	buf := new(bytes.Buffer)
	state := TexturedBlockState{
		Name:     b.tex.Name,
		Rotation: b.rotation,
	}
	if err := gob.NewEncoder(buf).Encode(state); err != nil {
		log.Panicf("texturedBlock.State() - %v", err)
	}
	return buf.Bytes()
}

func (b *texturedBlock) LoadState(state_bytes []byte) error {
	buf := bytes.NewBuffer(state_bytes)
	state := new(TexturedBlockState)
	if err := gob.NewDecoder(buf).Decode(state); err != nil {
		return err
	}
	b.tex = asset_loader.Texture(state.Name)
	b.rotation = state.Rotation
	return nil
}

// Dummy update method
func (b *compositeBlock) Update(world *World) {

}

func (b *compositeBlock) State() []byte {
	baseState := b.baseBlock.State()
	texState := b.texturedBlock.State()

	state := CompositeBlockState{
		BaseBlockState:     baseState,
		TexturedBlockState: texState,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(state); err != nil {
		log.Panicf("compositeBlock.State() - %v", err)
	}

	return buf.Bytes()
}

func (b *compositeBlock) LoadState(state_bytes []byte) error {
	buf := bytes.NewBuffer(state_bytes)
	state := new(CompositeBlockState)
	if err := gob.NewDecoder(buf).Decode(state); err != nil {
		return err
	}

	if err := b.baseBlock.LoadState(state.BaseBlockState); err != nil {
		return err
	}

	if err := b.texturedBlock.LoadState(state.TexturedBlockState); err != nil {
		return err
	}

	return nil
}

// Update connected sides
func (b *connectedBlock) Update(world *World) {

}

func (b *connectedBlock) Render(world *World, screen *ebiten.Image, pos util.Coords2f) {
	var sidesConnected [4]bool
	for i, side := range [4]util.Coords2i{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		x, y := int64(b.x)+side.X, int64(b.y)+side.Y
		neighbor, err := world.BlockAt(x, y)
		if err != nil {
			continue
		}

		if neighbor.Ground.Type() == b.Type() {
			sidesConnected[i] = true
			// If neighbor is on another chunk, trigger redraw of that chunk
			if neighbor.Ground.ParentChunk() != b.parentChunk {
				neighbor.Ground.ParentChunk().needsRedraw = true
			}
		}
	}
	b.tex.SetSidesArray(sidesConnected)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(asset_loader.Texture(b.tex.FullName()).Texture, opts)
}

func (b *connectedBlock) TextureName() string {
	return b.tex.FullName()
}

func (b *connectedBlock) State() []byte {
	baseState := b.baseBlock.State()

	texStateBuf := bytes.NewBuffer(make([]byte, 0))
	if err := gob.NewEncoder(texStateBuf).Encode(ConnectedTextureState{Base: b.tex.Base, SidesConnected: b.tex.SidesConnected}); err != nil {
		log.Panicf("connectedBlock.State() - %v", err)
	}
	texState := texStateBuf.Bytes()

	state := ConnectedBlockState{
		BaseBlockState:        baseState,
		ConnectedTextureState: texState,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(state); err != nil {
		log.Panicf("connectedBlock.State() - %v", err)
	}

	return buf.Bytes()
}

func (b *connectedBlock) LoadState(state_bytes []byte) error {
	buf := bytes.NewBuffer(state_bytes)
	state := new(ConnectedBlockState)
	if err := gob.NewDecoder(buf).Decode(state); err != nil {
		return err
	}

	if err := b.baseBlock.LoadState(state.BaseBlockState); err != nil {
		return err
	}

	tex := new(ConnectedTextureState)
	if err := gob.NewDecoder(bytes.NewReader(state.ConnectedTextureState)).Decode(tex); err != nil {
		return err
	}
	b.tex.Base = tex.Base
	b.tex.SidesConnected = tex.SidesConnected

	return nil
}

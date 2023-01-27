package types

import (
	"github.com/aquilax/go-perlin"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type Chunk interface {
	At(x uint, y uint) (Block, error)
	BlockCoords() Coords2u
	Coords() Coords2u
	Generate(baseGenerator *perlin.Perlin, secondaryGenerator *perlin.Perlin, mountainGenerator *perlin.Perlin) error
	GenerateDummy() error
	Render(world World)
	Save(id uuid.UUID) error
	SetBlock(x uint, y uint, block Block) error
	Update(world World)
	TriggerRedraw()
	Texture() *ebiten.Image
}

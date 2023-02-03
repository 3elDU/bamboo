package types

import (
	"github.com/aquilax/go-perlin"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type Chunk interface {
	// Returns a dummy block, in case of an error
	At(x uint, y uint) Block
	BlockCoords() Coords2u
	Coords() Coords2u
	Generate(baseGenerator *perlin.Perlin, secondaryGenerator *perlin.Perlin, mountainGenerator *perlin.Perlin)
	GenerateDummy()
	Render(world World)
	Save(id uuid.UUID)
	SetBlock(x uint, y uint, block Block)
	Update(world World)
	TriggerRedraw()
	Texture() *ebiten.Image
}

package world

import (
	"github.com/3elDU/bamboo/util"
)

type Block interface {
	Coords() util.Coords2i
	SetCoords(coords util.Coords2i)
	ParentChunk() *chunk
	SetParentChunk(chunk *chunk)

	Update()
	Render(target util.Coords2f)
}

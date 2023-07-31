package blocks_impl

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"

	"github.com/3elDU/bamboo/types"

	"github.com/3elDU/bamboo/assets"
)

// Max amount of ticks that one berry can take to grow
const BerryGrowthTime = 60 * 60 * 5

func init() {
	gob.Register(BerryBushState{})
	types.NewBerryBushBlock = NewBerryBushBlock
}

type BerryBushState struct {
	BaseBlockState

	DriedOut           bool
	Berries            int
	TotalBerriesGrown  int
	TicksTillNextBerry int
}

type BerryBushBlock struct {
	baseBlock
	texturedBlock

	// One bush can only grow a limited capacity of berries, after which it will "dry out"
	// To be able to grow berries again, bush needs watering, which can be done with the funnel.
	// Then, the process repeats.
	// This can be avoided by growing a bush directly near a water source.
	driedOut           bool
	berries            int
	totalBerriesGrown  int
	ticksTillNextBerry int
}

func NewBerryBushBlock(berries int) types.Block {
	return &BerryBushBlock{
		baseBlock: baseBlock{
			blockType: types.BerryBushBlock,
		},
		texturedBlock: texturedBlock{
			tex: assets.Texture(fmt.Sprintf("bush%v", berries)),
		},
		driedOut:           false,
		berries:            berries,
		totalBerriesGrown:  berries,
		ticksTillNextBerry: rand.Intn(BerryGrowthTime),
	}
}

func (b *BerryBushBlock) setBerries(berries int) {
	b.berries = berries
	if b.berries > 4 {
		b.berries = 4
	}
	b.texturedBlock = texturedBlock{
		tex: assets.Texture(fmt.Sprintf("bush%v", b.berries)),
	}
}

func (b *BerryBushBlock) Update(world types.World) {
	if b.ticksTillNextBerry <= 0 && !b.driedOut {
		b.setBerries(b.berries + 1)
		b.totalBerriesGrown += 1
		log.Printf("Total berries grown: %v", b.totalBerriesGrown)

		if b.berries < 4 {
			b.ticksTillNextBerry = rand.Intn(BerryGrowthTime)
		}

		b.parentChunk.MarkAsModified()
	}

	// Bush can grow 64 berries, after that it dries out
	if b.totalBerriesGrown >= 5 && !b.driedOut {
		b.driedOut = true
		b.parentChunk.MarkAsModified()
	}
	// Change the texture to dried out bush, when there are no more berries left
	if b.driedOut && b.berries == 0 {
		// Do not change the texture if it was already changed
		if b.tex.Name() != "dried_out_bush" {
			b.tex = assets.Texture("dried_out_bush")
			b.parentChunk.MarkAsModified()
		}
	}

	b.ticksTillNextBerry -= 1
}

func (b *BerryBushBlock) Break() {
	if b.berries == 0 {
		return
	}

	if types.GetInventory().AddItem(types.ItemSlot{
		Item:     types.NewBerryItem(),
		Quantity: 1,
	}) {
		b.setBerries(b.berries - 1)
		b.parentChunk.MarkAsModified()
	}
}

func (b *BerryBushBlock) State() interface{} {
	return BerryBushState{
		BaseBlockState:     b.baseBlock.State().(BaseBlockState),
		DriedOut:           b.driedOut,
		Berries:            b.berries,
		TotalBerriesGrown:  b.totalBerriesGrown,
		TicksTillNextBerry: b.ticksTillNextBerry,
	}
}

func (b *BerryBushBlock) LoadState(s interface{}) {
	state := s.(BerryBushState)
	b.baseBlock.LoadState(state.BaseBlockState)

	b.driedOut = state.DriedOut
	if b.driedOut {
		b.tex = assets.Texture("dried_out_bush")
	}

	b.setBerries(state.Berries)
	b.totalBerriesGrown = state.TotalBerriesGrown
	b.ticksTillNextBerry = state.TicksTillNextBerry
}

package inventory

import (
	"github.com/3elDU/bamboo/types"
	"github.com/google/uuid"
)

type SavedSlot struct {
	Quantity uint8
	ItemType types.ItemType
}

func LoadInventory(baseUUID uuid.UUID) *Inventory {
	// path := filepath.Join(config.WorldSaveDirectory, baseUUID.String(), config.InventoryFile)
	return nil
}

func (inv *Inventory) Save() {

}

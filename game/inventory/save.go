package inventory

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"path/filepath"

	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/google/uuid"
)

func LoadInventory(baseUUID uuid.UUID) *Inventory {
	path := filepath.Join(config.WorldSaveDirectory, baseUUID.String(), config.InventoryFile)

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read inventory save file: %v", err)
		return NewInventory()
	}

	loadedInventory := make([]types.SavedSlot, Size)
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&loadedInventory); err != nil {
		log.Printf("failed to decode inventory: %v", err)
		return NewInventory()
	}

	inventory := NewInventory()
	for i := 0; i < Size; i++ {
		if loadedInventory[i].Empty {
			continue
		}
		item := types.NewItem(loadedInventory[i].ItemType)
		item.LoadState(loadedInventory[i].State)
		inventory.Slots[i] = &types.ItemSlot{
			Empty:    false,
			Quantity: loadedInventory[i].Quantity,
			Item:     item,
		}
	}

	return inventory
}

func (inv *Inventory) Save(metadata types.Save) {
	path := filepath.Join(config.WorldSaveDirectory, metadata.BaseUUID.String(), config.InventoryFile)

	file, err := os.Create(path)
	if err != nil {
		log.Panicf("failed to open inventory file for saving: %v", err)
	}
	defer file.Close()

	savedInventory := make([]types.SavedSlot, Size)
	for i := 0; i < Size; i++ {
		if inv.Slots[i].Empty {
			savedInventory[i] = types.SavedSlot{
				Empty: true,
			}
		} else {
			savedInventory[i] = inv.Slots[i].Save()
		}
	}

	if err := gob.NewEncoder(file).Encode(savedInventory); err != nil {
		log.Panicf("failed to encode inventory: %v", err)
	}
}

package items_impl

import (
	"encoding/gob"

	"github.com/3elDU/bamboo/types"
)

func init() {
	gob.Register(BaseItemState{})
}

type BaseItemState struct {
	Type types.ItemType
}

type baseItem struct {
	id types.ItemType
}

func (i *baseItem) Type() types.ItemType {
	return i.id
}

func (i *baseItem) Stackable() bool {
	return true
}

func (i *baseItem) State() interface{} {
	return BaseItemState{
		Type: i.id,
	}
}

func (i *baseItem) LoadState(s interface{}) {
	state := s.(BaseItemState)
	i.id = state.Type
}

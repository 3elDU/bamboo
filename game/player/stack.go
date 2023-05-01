package player

import (
	"encoding/gob"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
)

// Stack maintains a LIFO stack of player states
// Top element on the stack is current state of the player.
// When player enters a new world, new player state is pushed onto the stack, and it becomes the current state
// When player goes back, we pop last element in the stack, going back to previous player state
type Stack struct {
	Stack []*Player
}

func NewPlayerStack() *Stack {
	return &Stack{
		Stack: make([]*Player, 0),
	}
}

func (stack *Stack) Top() *Player {
	last := len(stack.Stack) - 1
	if last < 0 {
		log.Panicf("Top() called on empty stack")
	}
	return stack.Stack[last]
}

func (stack *Stack) Push(player *Player) {
	stack.Stack = append(stack.Stack, player)
}

func (stack *Stack) Pop() *Player {
	top := stack.Top()
	stack.Stack = stack.Stack[:len(stack.Stack)-1]
	return top
}

func LoadPlayerStack(baseUUID uuid.UUID) *Stack {
	saveDir := filepath.Join(config.WorldSaveDirectory, baseUUID.String())

	f, err := os.Open(filepath.Join(saveDir, config.PlayerInfoFile))
	if err != nil {
		// if file does not exist, create a new stack
		return NewPlayerStack()
	}

	stack := new(Stack)
	if err := gob.NewDecoder(f).Decode(stack); err != nil {
		log.Panicf("LoadPlayerStack() - failed to decode metadata - %v", err)
	}

	return stack
}

func (stack *Stack) Save(metadata types.Save) {
	saveDir := filepath.Join(config.WorldSaveDirectory, metadata.BaseUUID.String())

	// make a save directory, if it doesn't exist yet
	os.Mkdir(saveDir, os.ModePerm)

	// open world metadata file
	f, err := os.Create(filepath.Join(saveDir, config.PlayerInfoFile))
	if err != nil {
		log.Panicf("failed to create player metadata file")
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(stack); err != nil {
		log.Panicf("failed to write player metadata")
	}
}

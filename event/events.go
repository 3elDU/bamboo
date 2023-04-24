package event

import "github.com/google/uuid"

// Enumeration with all declared event types
const (
	CaveEntered Type = iota
)

type CaveEnteredArgs struct {
	ID uuid.UUID
}

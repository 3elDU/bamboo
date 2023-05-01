package event

import "github.com/google/uuid"

// Enumeration with all declared event types
const (
	CaveEnter Type = iota
	CaveExit
)

type CaveEnteredArgs struct {
	ID uuid.UUID
}

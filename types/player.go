package types

var currentPlayer Player

func SetCurrentPlayer(player Player) {
	currentPlayer = player
}
func GetCurrentPlayer() Player {
	return currentPlayer
}

type Player interface {
	Position() Vec2f
	SetPosition(pos Vec2f)
	Move(delta Vec2f)
	Velocity() Vec2f

	LookingAt() Vec2u
}

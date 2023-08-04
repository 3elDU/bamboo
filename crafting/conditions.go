// Some crafts may be available only under certain conditions.
// E.g. a player must be near a campfire to craft an item from clay.

package crafting

import "github.com/3elDU/bamboo/types"

func PlayerMustBeNearCampfire() bool {
	position := types.GetCurrentPlayer().Position()
	const maxDistanceToCampfire = 4

	// Search for a campfire
	for x := position.X - maxDistanceToCampfire; x < position.X+maxDistanceToCampfire; x++ {
		for y := position.Y - maxDistanceToCampfire; y < position.Y+maxDistanceToCampfire; y++ {
			if campfire, ok := types.GetCurrentWorld().BlockAt(uint64(x), uint64(y)).(types.ICampfireBlock); ok {
				if campfire.IsLitUp() {
					return true
				}
			}
		}
	}

	return false
}

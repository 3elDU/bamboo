/*
	Things related to player rendering
*/

package player

import (
	"github.com/3elDU/bamboo/asset_loader"
	"github.com/hajimehoshi/ebiten/v2"
)

var texture_map = map[MovementDirection]string{
	Left:  "player_left",
	Right: "player_right",
	Up:    "player_up",
	Down:  "player_down",
	Still: "player_still",
}

func (player Player) Render(screen *ebiten.Image, scaling float64) {
	opts := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	tex := asset_loader.Texture(texture_map[player.movementDirection]).Texture()

	opts.GeoM.Reset()
	opts.GeoM.Scale(scaling, scaling)
	opts.GeoM.Translate(
		float64(sw)/2-8*scaling,
		float64(sh)/2-16*scaling,
	)
	screen.DrawImage(tex, opts)
}

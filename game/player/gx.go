/*
	Things related to player rendering
*/

package player

import (
	"image"
	"time"

	"github.com/3elDU/bamboo/asset_loader"
	"github.com/hajimehoshi/ebiten/v2"
)

var textureMap = map[MovementDirection]string{
	Left:  "player_left",
	Right: "player_right",
	Up:    "player_up",
	Down:  "player_down",
}

func (player *Player) Render(screen *ebiten.Image, scaling float64, paused bool) {
	opts := &ebiten.DrawImageOptions{}
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	tex := ebiten.NewImageFromImage(
		asset_loader.Texture(textureMap[player.MovementDirection]).Texture().SubImage(
			image.Rect(int(player.animationFrame)*16, 0, int(player.animationFrame)*16+16, 32),
		),
	)

	opts.GeoM.Reset()
	opts.GeoM.Scale(scaling, scaling)
	opts.GeoM.Translate(
		float64(sw)/2-8*scaling,
		float64(sh)/2-16*scaling,
	)
	screen.DrawImage(tex, opts)

	if !paused {
		player.nextAnimationFrame()
	}
}

func (player *Player) nextAnimationFrame() {
	// run at precisely 5 fps
	if time.Since(player.lastFrameChange).Seconds() < 0.2 {
		return
	}

	// if the player is standing still, reset the frame to 0
	if player.speed() < 0.01 {
		player.animationFrame = 0
		return
	}

	if player.animationFrame >= 3 {
		player.animationFrame = 0
	} else {
		player.animationFrame++
	}

	player.lastFrameChange = time.Now()
}

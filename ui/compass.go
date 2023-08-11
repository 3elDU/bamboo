package ui

import (
	"math"

	"github.com/3elDU/bamboo/assets"
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/types"
	"github.com/hajimehoshi/ebiten/v2"
)

type CompassComponent struct {
	baseComponent

	angleRad float64

	compassTexture *ebiten.Image
	arrowTexture   *ebiten.Image
	opts           *ebiten.DrawImageOptions
}

func NewCompassComponent() *CompassComponent {
	return &CompassComponent{
		baseComponent: newBaseComponent(),

		compassTexture: assets.Texture("compass").Texture(),
		arrowTexture:   assets.Texture("compass_arrow").Texture(),
		opts:           &ebiten.DrawImageOptions{},
	}
}

func (compass *CompassComponent) ComputedSize() (float64, float64) {
	size := compass.compassTexture.Bounds().Size()
	return float64(size.X) * config.UIScaling, float64(size.Y) * config.UIScaling
}
func (compass *CompassComponent) MaxSize() (float64, float64) {
	return compass.ComputedSize()
}
func (compass *CompassComponent) CapacityForChild(_ Component) (float64, float64) {
	return 0, 0
}
func (compass *CompassComponent) MaxCapacityForChild(_ Component) (float64, float64) {
	return 0, 0
}
func (compass *CompassComponent) Children() []Component {
	return []Component{}
}
func (compass *CompassComponent) Update() error {
	// Calculate the angle between the player's position and the world spawn point

	spawnpoint := types.GetCurrentWorld().PlayerSpawnPoint()
	playerPosition := types.GetCurrentPlayer().Position()
	deltaX := int(spawnpoint.X) - int(playerPosition.X)
	deltaY := int(spawnpoint.Y) - int(playerPosition.Y)
	compass.angleRad = math.Atan2(float64(deltaY), float64(deltaX))

	return nil
}
func (compass *CompassComponent) Draw(screen *ebiten.Image, x, y float64) error {
	compass.opts.GeoM.Reset()

	// Draw the compass background
	compass.opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	compass.opts.GeoM.Translate(x, y)
	screen.DrawImage(compass.compassTexture, compass.opts)

	// Draw the compass arrow
	compass.opts.GeoM.Reset()
	// Set the rotation origin to the center of the texture
	compass.opts.GeoM.Translate(-16, -16)
	compass.opts.GeoM.Rotate(compass.angleRad)
	compass.opts.GeoM.Translate(16, 16)
	compass.opts.GeoM.Scale(config.UIScaling, config.UIScaling)
	compass.opts.GeoM.Translate(x, y)
	screen.DrawImage(compass.arrowTexture, compass.opts)

	return nil
}

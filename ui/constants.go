package ui

import (
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
)

// Em equals letter height divided by 2
var Em = font.CharHeight / 2 * float64(config.UIScaling)

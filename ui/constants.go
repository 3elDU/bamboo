package ui

import (
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
)

// Em equals letter height divided by 2
var Em = font.CharHeight / 2 * float64(config.UIScaling)

// Minimum width for input elements, such as buttons and input fields
var MinInputWidth = float64(font.CharWidth) * 24 * config.UIScaling

package ui

import (
	"github.com/3elDU/bamboo/config"
	"github.com/3elDU/bamboo/font"
)

// Em equals letter height divided by 2
var Em float64 = font.FontHeight / 2 * float64(config.UIScaling)

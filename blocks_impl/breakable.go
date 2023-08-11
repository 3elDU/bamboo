package blocks_impl

import "github.com/3elDU/bamboo/types"

type breakableBlock struct {
	toolRequiredToBreak  types.ToolFamily
	toolStrengthRequired types.ToolStrength
}

func (b breakableBlock) ToolRequiredToBreak() types.ToolFamily {
	return b.toolRequiredToBreak
}
func (b breakableBlock) ToolStrengthRequired() types.ToolStrength {
	return b.toolStrengthRequired
}

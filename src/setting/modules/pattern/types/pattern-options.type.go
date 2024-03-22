package types

type PatternOptionsType struct {
	MaxStroke   float32 `json:"max_stroke" xml:"max_stroke"`
	MaxScale    int     `json:"max_scale" xml:"max_scale"`
	MaxSpacingX float32 `json:"max_spacing_x" xml:"max_spacing_x"`
	MaxSpacingY float32 `json:"max_spacing_y" xml:"max_spacing_y"`
}

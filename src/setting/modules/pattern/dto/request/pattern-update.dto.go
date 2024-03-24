package request

type PatternUpdateOptionDto struct {
	MaxStroke   float32 `json:"max_stroke" xml:"max_stroke" binding:"omitempty,gte=1"`
	MaxScale    int     `json:"max_scale" xml:"max_scale" binding:"omitempty,gte=1"`
	MaxSpacingX float32 `json:"max_spacing_x" xml:"max_spacing_x" binding:"omitempty,gte=1"`
	MaxSpacingY float32 `json:"max_spacing_y" xml:"max_spacing_y" binding:"omitempty,gte=1"`
}

type PatternUpdateDto struct {
	Mode   string  `json:"mode" xml:"mode" binding:"omitempty,oneof=fill stroke join"`
	Colors int     `json:"colors" xml:"colors" binding:"omitempty,gte=1"`
	Width  float32 `json:"width" xml:"width" binding:"omitempty,gte=1"`
	Height float32 `json:"height" xml:"height" binding:"omitempty,gte=1"`
	Path   string  `json:"path" xml:"path" binding:"omitempty,html,startswith=<path,endswith=/>"`

	Options *PatternUpdateOptionDto `json:"options" xml:"options" binding:"omitempty"`
}

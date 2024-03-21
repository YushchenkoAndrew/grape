package request

type PatternCreateOptionDto struct {
	MaxStroke   float32 `json:"max_stroke" xml:"max_stroke" binding:"required,gte=1"`
	MaxScale    int     `json:"max_scale" xml:"max_scale" binding:"required,gte=1"`
	MaxSpacingX float32 `json:"max_spacing_x" xml:"max_spacing_x" binding:"required,gte=1"`
	MaxSpacingY float32 `json:"max_spacing_y" xml:"max_spacing_y" binding:"required,gte=1"`
}

type PatternCreateDto struct {
	Mode   string  `json:"mode" xml:"mode" binding:"required,oneof=fill stroke join"`
	Colors int     `json:"colors" xml:"colors" binding:"required,gte=1"`
	Width  float32 `json:"width" xml:"width" binding:"required,gte=1"`
	Height float32 `json:"height" xml:"height" binding:"required,gte=1"`
	Path   string  `json:"path" xml:"path" binding:"required,html,startswith=<path,endswith=/>"`

	Options *PatternCreateOptionDto `json:"options" xml:"options" binding:"required"`
}

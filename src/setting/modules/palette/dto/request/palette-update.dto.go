package request

type PaletteUpdateDto struct {
	Mode   string  `json:"mode" xml:"mode" binding:"omitempty,oneof=fill stroke join"`
	Colors int     `json:"colors" xml:"colors" binding:"omitempty,gte=1"`
	Width  float32 `json:"width" xml:"width" binding:"omitempty,gte=1"`
	Height float32 `json:"height" xml:"height" binding:"omitempty,gte=1"`
	Path   string  `json:"path" xml:"path" binding:"omitempty,html,startswith=<path,endswith=/>"`
}

package response

type SvgPatternBasicResponseDto struct {
	Mode   string  `copier:"GetMode" json:"mode" xml:"mode" example:"fill"`
	Width  float32 `json:"width" xml:"width" example:"10"`
	Height float32 `json:"height" xml:"heigh" example:"heigh"`
	Path   string  `json:"path" xml:"path" example:"<path src='' />"`
}

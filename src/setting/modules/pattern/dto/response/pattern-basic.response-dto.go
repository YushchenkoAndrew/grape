package response

type PatternBasicResponseDto struct {
	Id     string  `copier:"UUID" json:"id" xml:"id" example:"a3c94c88-944d-4636-86d1-c233bdfaf70e"`
	Mode   string  `copier:"GetMode" json:"mode" xml:"mode" example:"fill"`
	Width  float32 `json:"width" xml:"width" example:"10"`
	Height float32 `json:"height" xml:"heigh" example:"heigh"`
	Path   string  `json:"path" xml:"path" example:"<path src='' />"`
}

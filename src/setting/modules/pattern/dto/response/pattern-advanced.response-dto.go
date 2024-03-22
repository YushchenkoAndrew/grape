package response

type PatternAdvancedResponseDto struct {
	PatternBasicResponseDto
	Colors  int         `json:"colors" xml:"colors" example:"10"`
	Options interface{} `copier:"GetOptions" json:"options" xml:"options"`
}

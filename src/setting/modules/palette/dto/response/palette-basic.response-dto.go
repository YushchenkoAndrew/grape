package response

type PaletteBasicResponseDto struct {
	Id     string   `copier:"UUID" json:"id" xml:"id" example:"a3c94c88-944d-4636-86d1-c233bdfaf70e"`
	Colors []string `json:"colors" xml:"colors"`
}

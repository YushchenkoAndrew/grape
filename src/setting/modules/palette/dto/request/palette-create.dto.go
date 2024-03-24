package request

type PaletteCreateDto struct {
	Colors []string `json:"colors" xml:"colors" binding:"required,dive,hexcolor"`
}

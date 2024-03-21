package request

type PaletteCreateDto struct {
	Colors []string `json:"colors" xml:"colors" binding:"required,regexp=^#?([a-f0-9]{6}|[a-f0-9]{3})$"`
}

package request

type LinkUpdateDto struct {
	Name string `json:"name" xml:"name" binding:"omitempty"`
	Link string `json:"link" xml:"link" binding:"omitempty,url"`
}

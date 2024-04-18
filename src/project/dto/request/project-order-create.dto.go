package request

type ProjectOrderUpdateDto struct {
	Position int `json:"position" xml:"position" binding:"required,gte=1"`
}

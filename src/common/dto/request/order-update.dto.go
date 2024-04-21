package request

type OrderUpdateDto struct {
	Position int `json:"position" xml:"position" binding:"required,gte=1"`
}

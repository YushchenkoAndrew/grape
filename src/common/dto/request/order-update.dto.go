package request

type OrderUpdateDto struct {
	Position int `json:"position" xml:"position" binding:"omitempty,gte=0"`
}

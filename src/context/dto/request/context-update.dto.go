package request

type ContextUpdateDto struct {
	Name string `json:"name" xml:"name" binding:"omitempty"`
}

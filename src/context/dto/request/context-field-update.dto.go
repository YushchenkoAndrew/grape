package request

type ContextFieldUpdateDto struct {
	Name  string `json:"name" xml:"name" binding:"omitempty"`
	Value string `json:"value" xml:"value" binding:"omitempty"`
}

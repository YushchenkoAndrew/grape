package request

type TagUpdateDto struct {
	Name string `json:"name" xml:"name" binding:"omitempty"`
}

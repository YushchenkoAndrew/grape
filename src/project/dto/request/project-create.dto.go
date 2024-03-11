package request

type ProjectCreateDto struct {
	Name        string `json:"name" xml:"name" binding:"required"`
	Description string `json:"description" xml:"description" binding:"omitempty"`
	Type        string `json:"type" xml:"type" binding:"required,oneof=html markdown link k3s"`
	Footer      string `json:"footer" xml:"footer" binding:"omitempty"`
}

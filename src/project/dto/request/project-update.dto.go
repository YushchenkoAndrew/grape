package request

type ProjectUpdateDto struct {
	Name        string `json:"name" xml:"name" binding:"omitempty"`
	Description string `json:"description" xml:"description" binding:"omitempty"`
	// Type        string `json:"type" xml:"type" binding:"omitempty,oneof=html markdown link k3s"`
	Status string `json:"status" xml:"status" binding:"omitempty,oneof=active inactive"`
	Footer string `json:"footer" xml:"footer" binding:"omitempty"`

	// Attachments *
}

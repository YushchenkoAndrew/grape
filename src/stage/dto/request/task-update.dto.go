package request

type TaskUpdateDto struct {
	Name        string `json:"name" xml:"name" binding:"omitempty"`
	Description string `json:"description" xml:"description" binding:"omitempty"`
	Status      string `json:"status" xml:"status" binding:"omitempty,oneof=active inactive"`
}

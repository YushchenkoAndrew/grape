package request

type StageUpdateDto struct {
	Name   string `json:"name" xml:"name" binding:"omitempty"`
	Status string `json:"status" xml:"status" binding:"omitempty,oneof=active inactive"`
}

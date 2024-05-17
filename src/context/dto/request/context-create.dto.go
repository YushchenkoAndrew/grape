package request

type ContextCreateDto struct {
	Name            string `json:"name" xml:"name" binding:"required"`
	ContextableID   string `json:"contextable_id" xml:"contextable_id" binding:"required,uuid4"`
	ContextableType string `json:"contextable_type" xml:"contextable_type" binding:"required,oneof=tasks"`
}

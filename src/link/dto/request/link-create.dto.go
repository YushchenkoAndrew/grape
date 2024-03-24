package request

type LinkCreateDto struct {
	Name         string `json:"name" xml:"name" binding:"required"`
	Link         string `json:"link" xml:"link" binding:"required,url"`
	LinkableID   string `json:"linkable_id" xml:"linkable_id" binding:"required,uuid4"`
	LinkableType string `json:"linkable_type" xml:"linkable_type" binding:"required,oneof=projects"`
}

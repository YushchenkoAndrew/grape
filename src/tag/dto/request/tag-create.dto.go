package request

type TagCreateDto struct {
	Name         string `json:"name" xml:"name" binding:"required"`
	TaggableID   string `json:"taggable_id" xml:"taggable_id" binding:"required,uuid4"`
	TaggableType string `json:"taggable_type" xml:"taggable_type" binding:"required,oneof=projects tasks"`
}

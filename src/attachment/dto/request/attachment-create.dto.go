package request

type AttachmentCreateDto struct {
	Path           string `form:"path" binding:"required,startswith=/,dirpath"`
	AttachableID   string `form:"attachable_id" binding:"required,uuid4"`
	AttachableType string `form:"attachable_type" binding:"required,oneof=projects"`
}

package request

type AttachmentCreateDto struct {
	Path           string `form:"path" binding:"required,startswith=/"`
	AttachableID   string `form:"attachable_id" binding:"required"`
	AttachableType string `form:"attachable_type" binding:"required,oneof=projects"`
}

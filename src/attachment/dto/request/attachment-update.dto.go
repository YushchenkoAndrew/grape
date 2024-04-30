package request

type AttachmentUpdateDto struct {
	Name    string `form:"name" binding:"omitempty,filepath"`
	Path    string `form:"path" binding:"omitempty,startswith=/,dirpath"`
	Preview bool   `form:"preview" binding:"omitempty"`
}

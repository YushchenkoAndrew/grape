package response

type AttachmentAdvancedResponseDto struct {
	AttachmentBasicResponseDto

	CreatedAt string `copier:"CreatedAtISO" json:"created_at" xml:"created_at"`
	Size      int64  `json:"size" xml:"size" example:"100"`
}

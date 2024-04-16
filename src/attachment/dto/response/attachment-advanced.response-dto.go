package response

type AttachmentAdvancedResponseDto struct {
	AttachmentBasicResponseDto

	UpdatedAt string `copier:"UpdatedAtISO" json:"updated_at" xml:"updated_at"`
	Size      int64  `json:"size" xml:"size" example:"100"`
}

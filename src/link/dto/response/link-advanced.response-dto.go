package response

type LinkAdvancedResponseDto struct {
	LinkBasicResponseDto

	UpdatedAt string `copier:"UpdatedAtISO" json:"updated_at" xml:"updated_at"`
	Size      int64  `json:"size" xml:"size" example:"100"`
	Order     int    `json:"order" xml:"order" example:"100"`
}

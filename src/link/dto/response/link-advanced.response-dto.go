package response

type LinkAdvancedResponseDto struct {
	LinkBasicResponseDto

	UpdatedAt string `copier:"UpdatedAtISO" json:"updated_at" xml:"updated_at"`
	Order     int    `json:"order" xml:"order" example:"100"`
}

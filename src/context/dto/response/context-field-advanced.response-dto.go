package response

type ContextFieldAdvancedResponseDto struct {
	ContextFieldBasicResponseDto

	UpdatedAt string `copier:"UpdatedAtISO" json:"updated_at" xml:"updated_at"`
	Order     int    `json:"order" xml:"order" example:"100"`
}

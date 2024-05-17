package response

type ContextAdvancedResponseDto struct {
	ContextBasicResponseDto

	UpdatedAt string `copier:"UpdatedAtISO" json:"updated_at" xml:"updated_at"`
	Order     int    `json:"order" xml:"order" example:"100"`

	Fields []ContextFieldAdvancedResponseDto `copier:"ContextFields" json:"fields" xml:"fields"`
}

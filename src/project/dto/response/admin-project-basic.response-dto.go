package response

import (
	att "grape/src/attachment/dto/response"
	ln "grape/src/link/dto/response"
	user "grape/src/user/dto/response"
)

type AdminProjectBasicResponseDto struct {
	ProjectBasicResponseDto

	CreatedAt string `copier:"CreatedAtISO" json:"created_at" xml:"created_at"`
	Order     int    `json:"order" xml:"order" example:"0"`
	Status    string `copier:"GetStatus" json:"status" xml:"status" example:"true"`

	Owner       user.UserBasicResponseDto           `json:"owner" xml:"owner"`
	Attachments []att.AttachmentAdvancedResponseDto `json:"attachments" xml:"attachments"`
	Links       []ln.LinkBasicResponseDto           `json:"links" xml:"links"`
}

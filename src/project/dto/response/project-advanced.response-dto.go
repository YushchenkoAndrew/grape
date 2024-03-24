package response

import (
	user "grape/src/user/dto/response"
)

type ProjectAdvancedResponseDto struct {
	ProjectDetailedResponseDto

	CreatedAt string `copier:"CreatedAtISO" json:"created_at" xml:"created_at"`
	Order     int    `json:"order" xml:"order" example:"0"`
	Status    string `copier:"GetStatus" json:"status" xml:"status" example:"active"`

	Owner user.UserBasicResponseDto `json:"owner" xml:"owner"`
	// Attachments []att.AttachmentAdvancedResponseDto `json:"attachments" xml:"attachments"`
}

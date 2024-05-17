package response

import (
	att "grape/src/attachment/dto/response"
	ctx "grape/src/context/dto/response"
	ln "grape/src/link/dto/response"
)

type AdminTaskBasicResponseDto struct {
	TaskBasicResponseDto

	Order  int    `json:"order" xml:"order" example:"0"`
	Status string `copier:"GetStatus" json:"status" xml:"status" example:"true"`

	Links       []ln.LinkAdvancedResponseDto        `json:"links" xml:"links"`
	Contexts    []ctx.ContextAdvancedResponseDto    `json:"contexts" xml:"contexts"`
	Attachments []att.AttachmentAdvancedResponseDto `json:"attachments" xml:"attachments"`
}

package response

import (
	att "grape/src/attachment/dto/response"
	r "grape/src/common/dto/response"
	ctx "grape/src/context/dto/response"
	ln "grape/src/link/dto/response"
	user "grape/src/user/dto/response"
)

type TaskBasicResponseDto struct {
	r.UuidResponseDto

	Description string `json:"description" xml:"description" example:"Description"`

	Owner       user.UserBasicResponseDto        `json:"owner" xml:"owner"`
	Attachments []att.AttachmentBasicResponseDto `json:"attachments" xml:"attachments"`
	Links       []ln.LinkBasicResponseDto        `json:"links" xml:"links"`
	Contexts    []ctx.ContextBasicResponseDto    `json:"contexts" xml:"contexts"`
}

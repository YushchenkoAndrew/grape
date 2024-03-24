package response

import (
	att "grape/src/attachment/dto/response"
	ln "grape/src/link/dto/response"
)

type ProjectDetailedResponseDto struct {
	ProjectBasicResponseDto

	Footer      string                           `json:"footer" xml:"footer" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`
	Attachments []att.AttachmentBasicResponseDto `json:"attachments" xml:"attachments"`
	Links       []ln.LinkBasicResponseDto        `json:"links" xml:"links"`
}

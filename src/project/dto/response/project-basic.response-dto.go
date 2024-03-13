package response

import (
	att "grape/src/attachment/dto/response"
	r "grape/src/common/dto/response"
)

type ProjectBasicResponseDto struct {
	r.UuidResponseDto

	Description string `json:"description" xml:"description" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
	Type        string `copier:"GetType" json:"type" xml:"type" example:"html"`
	Footer      string `json:"footer" xml:"footer" example:"Creating a 'Code Rain' effect from Matrix. As funny joke you can put any text to display at the end."`

	Attachments []att.AttachmentBasicResponseDto `json:"attachments" xml:"attachments"`
}

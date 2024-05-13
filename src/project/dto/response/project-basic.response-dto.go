package response

import (
	att "grape/src/attachment/dto/response"
	r "grape/src/common/dto/response"
	ln "grape/src/link/dto/response"
)

type ProjectBasicResponseDto struct {
	r.UuidResponseDto

	Description string `json:"description" xml:"description" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
	Type        string `copier:"GetType" json:"type" xml:"type" example:"html"`

	Redirect  *ln.LinkBasicResponseDto        `json:"redirect" xml:"redirect"`
	Thumbnail *att.AttachmentBasicResponseDto `json:"thumbnail" xml:"thumbnail"`
}

package response

import (
	att "grape/src/attachment/dto/response"
	r "grape/src/common/dto/response"
	style "grape/src/style/dto/response"
)

type ProjectBasicResponseDto struct {
	r.UuidResponseDto

	Description string   `json:"description" xml:"description" example:"Take the blue pill and the sit will close, or take the red pill and I show how deep the rebbit hole goes"`
	Type        string   `copier:"GetType" json:"type" xml:"type" example:"html"`
	Colors      []string `copier:"GetColors" json:"colors" xml:"colors"`

	Pattern   style.SvgPatternBasicResponseDto `copier:"SvgPattern" json:"pattern" xml:"pattern"`
	Thumbnail *att.AttachmentBasicResponseDto  `copier:"GetThumbnail" json:"thumbnail" xml:"thumbnail"`
}

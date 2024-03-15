package response

import r "grape/src/common/dto/response"

type AttachmentBasicResponseDto struct {
	r.UuidResponseDto

	Type string `json:"type" xml:"type" example:".png"`
	Path string `copier:"GetAttachment" json:"path" xml:"path" example:"/test"`
}

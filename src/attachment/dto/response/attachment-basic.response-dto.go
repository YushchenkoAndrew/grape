package response

import r "grape/src/common/dto/response"

type AttachmentBasicResponseDto struct {
	r.UuidResponseDto

	Type string `json:"type" xml:"type" example:".png"`
	Path string `json:"path" xml:"path" example:"/test"`
	File string `copier:"GetAttachment" json:"file" xml:"file" example:"http://localhost/test.jpeg"`
}

package response

import r "grape/src/common/dto/response"

type ContextFieldBasicResponseDto struct {
	r.UuidResponseDto

	Value   string      `json:"value" xml:"value" example:"root"`
	Options interface{} `copier:"GetOptions" json:"options" xml:"options"`
}

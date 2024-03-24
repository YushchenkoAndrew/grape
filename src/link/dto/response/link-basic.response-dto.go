package response

import r "grape/src/common/dto/response"

type LinkBasicResponseDto struct {
	r.UuidResponseDto

	Name string `json:"name" xml:"name" example:"test"`
	Link string `json:"link" xml:"link" example:"http://test.com"`
}

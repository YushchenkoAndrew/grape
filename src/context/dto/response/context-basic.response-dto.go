package response

import r "grape/src/common/dto/response"

type ContextBasicResponseDto struct {
	r.UuidResponseDto

	Fields []ContextFieldBasicResponseDto `json:"fields" xml:"fields"`
}

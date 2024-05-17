package response

import (
	r "grape/src/common/dto/response"
)

type StageBasicResponseDto struct {
	r.UuidResponseDto

	Tasks []TaskBasicResponseDto `json:"tasks" xml:"tasks"`
}

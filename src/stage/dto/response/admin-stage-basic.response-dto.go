package response

import (
	r "grape/src/common/dto/response"
)

type AdminStageBasicResponseDto struct {
	r.UuidResponseDto

	Order  int    `json:"order" xml:"order" example:"0"`
	Status string `copier:"GetStatus" json:"status" xml:"status" example:"true"`

	Tasks []AdminTaskBasicResponseDto `json:"tasks" xml:"tasks"`
}

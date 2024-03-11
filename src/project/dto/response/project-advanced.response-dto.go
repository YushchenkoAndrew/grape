package response

import (
	"grape/src/user/dto/response"
)

type ProjectAdvancedResponseDto struct {
	ProjectBasicResponseDto

	Order  int                           `json:"order" xml:"order" example:"0"`
	Status bool                          `copier:"GetStatus" json:"status" xml:"status" example:"active"`
	Owner  response.UserBasicResponseDto `copier:"Owner" json:"owner" xml:"owner"`
}

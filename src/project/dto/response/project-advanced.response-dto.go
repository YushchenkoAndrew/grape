package response

import (
	"grape/src/user/dto/response"
)

type ProjectAdvancedResponseDto struct {
	ProjectBasicResponseDto

	Status bool                     `copier:"GetStatus" json:"status" xml:"status" example:"active"`
	Owner  response.UserResponseDto `json:"owner" xml:"owner"`
}

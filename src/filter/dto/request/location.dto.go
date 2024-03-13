package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type LocationDto struct {
	*request.CurrentUserDto

	IP []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func NewLocationDto(user *entities.UserEntity, init ...*LocationDto) *LocationDto {
	return request.NewRequest(&LocationDto{CurrentUserDto: request.NewCurrentUserDto(user)}, init...)
}

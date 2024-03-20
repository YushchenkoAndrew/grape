package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type StatisticDto struct {
	*request.CurrentUserDto

	ProjectIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func NewStatisticDto(user *entities.UserEntity, init ...*StatisticDto) *StatisticDto {
	return request.NewRequest(&StatisticDto{CurrentUserDto: request.NewCurrentUserDto(user)}, init...)
}

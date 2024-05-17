package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type LinkDto struct {
	*request.CurrentUserDto

	LinkIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *LinkDto) GetIds() []string {
	return c.LinkIds
}

func NewLinkDto(user *entities.UserEntity, init ...*LinkDto) *LinkDto {
	return request.NewRequest(&LinkDto{CurrentUserDto: request.NewCurrentUserDto(user)}, init...)
}

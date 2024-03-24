package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type LinkDto struct {
	*request.CurrentUserDto

	LinkIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *LinkDto) UUID() string {
	return c.LinkIds[0]
}

func NewLinkDto(user *entities.UserEntity, init ...*LinkDto) *LinkDto {
	return request.NewRequest(&LinkDto{CurrentUserDto: request.NewCurrentUserDto(user)}, init...)
}

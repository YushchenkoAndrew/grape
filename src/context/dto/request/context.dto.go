package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type ContextDto struct {
	*request.CurrentUserDto

	ContextIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *ContextDto) GetIds() []string {
	return c.ContextIds
}

func NewContextDto(user *entities.UserEntity, init ...*ContextDto) *ContextDto {
	return request.NewRequest(&ContextDto{CurrentUserDto: request.NewCurrentUserDto(user)}, init...)
}

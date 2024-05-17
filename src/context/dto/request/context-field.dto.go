package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type ContextFieldDto struct {
	*request.CurrentUserDto

	ContextIds      []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	ContextFieldIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *ContextFieldDto) GetIds() []string {
	return c.ContextFieldIds
}

func NewContextFieldDto(user *entities.UserEntity, init ...*ContextFieldDto) *ContextFieldDto {
	return request.NewRequest(&ContextFieldDto{CurrentUserDto: request.NewCurrentUserDto(user)}, init...)
}

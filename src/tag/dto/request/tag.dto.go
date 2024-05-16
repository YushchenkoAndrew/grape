package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type TagDto struct {
	*request.CurrentUserDto

	TagIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *TagDto) GetIds() []string {
	return c.TagIds
}

func NewTagDto(user *entities.UserEntity, init ...*TagDto) *TagDto {
	return request.NewRequest(&TagDto{CurrentUserDto: request.NewCurrentUserDto(user)}, init...)
}

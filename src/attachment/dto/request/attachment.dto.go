package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type AttachmentDto struct {
	*request.CurrentUserDto

	AttachmentIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *AttachmentDto) UUID() string {
	return c.AttachmentIds[0]
}

func NewAttachmentDto(user *entities.UserEntity, init ...*AttachmentDto) *AttachmentDto {
	return request.NewRequest(&AttachmentDto{CurrentUserDto: request.NewCurrentUserDto(user)}, init...)
}

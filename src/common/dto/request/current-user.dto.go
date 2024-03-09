package request

import (
	"grape/src/user/entities"

	"github.com/jinzhu/copier"
)

type CurrentUserDto struct {
	CurrentUser *entities.UserEntity
}

func (*CurrentUserDto) Offset() int {
	return 0
}

func (*CurrentUserDto) Limit() int {
	return 0
}

func NewCurrentUserDto(user *entities.UserEntity) *CurrentUserDto {
	return &CurrentUserDto{CurrentUser: user}
}

func NewRequest[Dto any](dst Dto, src ...Dto) Dto {
	for _, init := range src {
		copier.CopyWithOption(&dst, init, copier.Option{IgnoreEmpty: true})
	}
	return dst
}

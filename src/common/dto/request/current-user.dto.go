package request

import "grape/src/user/entities"

type CurrentUserDto struct {
	CurrentUser *entities.UserEntity
}

func NewCurrentUserDto(user *entities.UserEntity) CurrentUserDto {
	return CurrentUserDto{CurrentUser: user}
}

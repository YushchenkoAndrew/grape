package request

import "grape/src/user/entities"

type PageDto struct {
	*CurrentUserDto

	Page      int    `form:"page,omitempty,gte=1,default=1" example:"1"`
	Take      int    `form:"take,omitempty,gte=1,default=30" example:"20"`
	SortBy    string `form:"sort_by,omitempty,default=created_at" example:"name"`
	Direction string `form:"direction,omitempty,default=desc,oneof=asc desc" example:"asc"`
}

func (c *PageDto) Offset() int {
	return (c.Page - 1) * c.Take
}

func (c *PageDto) Limit() int {
	return c.Take
}

func NewPageDto(user *entities.UserEntity) *PageDto {
	return &PageDto{CurrentUserDto: NewCurrentUserDto(user)}
}

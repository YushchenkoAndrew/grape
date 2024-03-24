package request

import "grape/src/user/entities"

type PageDto struct {
	*CurrentUserDto

	Page      int    `form:"page,default=1" binding:"gte=1" example:"1"`
	Take      int    `form:"take,default=30" binding:"gte=1" example:"20"`
	SortBy    string `form:"sort_by,default=created_at" example:"name"`
	Direction string `form:"direction,default=desc" binding:"oneof=asc desc" example:"asc"`
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

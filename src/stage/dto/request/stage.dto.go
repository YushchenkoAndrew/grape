package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type StageDto struct {
	*request.PageDto

	SortBy    string `form:"sort_by,default=order" example:"name"`
	Direction string `form:"direction,default=asc" binding:"oneof=asc desc" example:"asc"`

	StageIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	Statuses []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *StageDto) UUID() string {
	return c.StageIds[0]
}

func NewStageDto(user *entities.UserEntity, init ...*StageDto) *StageDto {
	return request.NewRequest(&StageDto{PageDto: request.NewPageDto(user)}, init...)
}

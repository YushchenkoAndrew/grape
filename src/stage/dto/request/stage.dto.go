package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type StageDto struct {
	*request.PageDto

	Query     string `form:"query,omitempty" binding:"startsnotwith=%,endsnotwith=%" example:"test"`
	SortBy    string `form:"sort_by,default=order" example:"name"`
	Direction string `form:"direction,default=asc" binding:"oneof=asc desc" example:"asc"`

	StageIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	Statuses []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *StageDto) GetIds() []string {
	return c.StageIds
}

func NewStageDto(user *entities.UserEntity, init ...*StageDto) *StageDto {
	return request.NewRequest(&StageDto{PageDto: request.NewPageDto(user)}, init...)
}

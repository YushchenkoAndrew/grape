package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type TaskDto struct {
	*request.PageDto

	SortBy    string `form:"sort_by,default=order" example:"name"`
	Direction string `form:"direction,default=asc" binding:"oneof=asc desc" example:"asc"`

	TaskIds  []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	StageIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	Statuses []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *TaskDto) GetIds() []string {
	return c.TaskIds
}

func NewTaskDto(user *entities.UserEntity, init ...*TaskDto) *TaskDto {
	return request.NewRequest(&TaskDto{PageDto: request.NewPageDto(user)}, init...)
}

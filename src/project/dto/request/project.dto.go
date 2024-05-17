package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type ProjectDto struct {
	*request.PageDto

	Query     string `form:"query,omitempty" binding:"startsnotwith=%,endsnotwith=%" example:"test"`
	SortBy    string `form:"sort_by,default=order" example:"name"`
	Direction string `form:"direction,default=asc" binding:"oneof=asc desc" example:"asc"`

	ProjectIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	Statuses   []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	Types      []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *ProjectDto) GetIds() []string {
	return c.ProjectIds
}

func NewProjectDto(user *entities.UserEntity, init ...*ProjectDto) *ProjectDto {
	return request.NewRequest(&ProjectDto{PageDto: request.NewPageDto(user)}, init...)
}

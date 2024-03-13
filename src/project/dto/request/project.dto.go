package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type ProjectDto struct {
	*request.PageDto

	Query string `form:"query,omitempty" example:"test"`

	ProjectIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	Statuses   []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
	Types      []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *ProjectDto) UUID() string {
	return c.ProjectIds[0]
}

func NewProjectDto(user *entities.UserEntity, init ...*ProjectDto) *ProjectDto {
	return request.NewRequest(&ProjectDto{PageDto: request.NewPageDto(user)}, init...)
}

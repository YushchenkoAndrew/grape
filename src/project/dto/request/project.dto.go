package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type ProjectDto struct {
	*request.PageDto

	Query string `form:"query,omitempty" example:"test"`

	ProjectIds []string
	Statuses   []string
	Types      []string
}

func NewProjectDto(user *entities.UserEntity, init ...*ProjectDto) *ProjectDto {
	return request.NewRequest(&ProjectDto{PageDto: request.NewPageDto(user)}, init...)
}

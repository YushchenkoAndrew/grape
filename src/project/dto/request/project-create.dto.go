package request

import (
	"grape/src/common/dto/request"
	"grape/src/project/types"
	"grape/src/user/entities"
)

type ProjectCreateDto struct {
	Name        string `json:"name" xml:"name" binding:"required"`
	Description string `json:"description" xml:"description"`
	Type        string `json:"type" xml:"type" binding:"required,oneof=html markdown link k3s"`
	Footer      string `json:"footer" xml:"footer"`
}

func (c *ProjectCreateDto) GetType() types.ProjectTypeEnum {
	return types.Html.Value(c.Type)
}

func NewProjectCreateDto(user *entities.UserEntity, init ...*ProjectCreateDto) *ProjectCreateDto {
	return request.NewRequest(&ProjectCreateDto{}, init...)
}

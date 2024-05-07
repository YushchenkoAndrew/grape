package request

import "grape/src/stage/entities"

type TaskCreateDto struct {
	Name string `json:"name" xml:"name" binding:"required"`

	Stage *entities.StageEntity `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

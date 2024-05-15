package request

import "grape/src/context/entities"

type ContextFieldCreateDto struct {
	Name    string      `json:"name" xml:"name" binding:"required"`
	Options interface{} `json:"options" xml:"options" binding:"omitempty"`

	Context *entities.ContextEntity `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

package request

import "grape/src/context/entities"

type ContextFieldCreateDto struct {
	Name    string                  `json:"name" xml:"name" binding:"required"`
	Value   *string                 `json:"value" xml:"value" binding:"omitempty"`
	Options *map[string]interface{} `json:"options" xml:"options" binding:"omitempty"`

	Context *entities.ContextEntity `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

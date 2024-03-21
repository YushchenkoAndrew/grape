package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type PatternDto struct {
	*request.PageDto

	Modes  []string `form:"mode,omitempty,oneof=fill stroke join"`
	Colors []int    `form:"colors,omitempty"`

	PatternIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *PatternDto) UUID() string {
	return c.PatternIds[0]
}

func NewPatternDto(user *entities.UserEntity, init ...*PatternDto) *PatternDto {
	return request.NewRequest(&PatternDto{PageDto: request.NewPageDto(user)}, init...)
}

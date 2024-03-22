package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type PaletteDto struct {
	*request.PageDto

	// TODO: ?????
	// Colors []string `form:"colors,omitempty,regexp=^#?([a-f0-9]{6}|[a-f0-9]{3})$"`

	PaletteIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *PaletteDto) UUID() string {
	return c.PaletteIds[0]
}

func NewPaletteDto(user *entities.UserEntity, init ...*PaletteDto) *PaletteDto {
	return request.NewRequest(&PaletteDto{PageDto: request.NewPageDto(user)}, init...)
}

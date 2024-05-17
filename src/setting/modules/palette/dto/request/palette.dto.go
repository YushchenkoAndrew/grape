package request

import (
	"grape/src/common/dto/request"
	"grape/src/user/entities"
)

type PaletteDto struct {
	*request.PageDto

	PaletteIds []string `form:"-" json:"-" xml:"-" swaggerignore:"true"`
}

func (c *PaletteDto) GetIds() []string {
	return c.PaletteIds
}

func NewPaletteDto(user *entities.UserEntity, init ...*PaletteDto) *PaletteDto {
	return request.NewRequest(&PaletteDto{PageDto: request.NewPageDto(user)}, init...)
}

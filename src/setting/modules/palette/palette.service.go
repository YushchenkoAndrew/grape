package palette

import (
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/setting/modules/palette/dto/request"
	"grape/src/setting/modules/palette/dto/response"
	repo "grape/src/setting/modules/palette/repositories"
)

type PaletteService struct {
	Repository *repo.PaletteRepositoryT
}

func NewPaletteService(s *service.CommonService) *PaletteService {
	return &PaletteService{
		Repository: repo.NewPaletteRepository(s.DB),
	}
}

func (c *PaletteService) FindOne(dto *request.PaletteDto) (*response.PaletteBasicResponseDto, error) {
	palette, err := c.Repository.ValidateEntityExistence(dto)

	return common.NewResponse[response.PaletteBasicResponseDto](palette), err
}

func (c *PaletteService) FindAll(dto *request.PaletteDto) (*common.PageResponseDto[[]response.PaletteBasicResponseDto], error) {
	total, palette, err := c.Repository.GetAllPage(dto)

	return &common.PageResponseDto[[]response.PaletteBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.PaletteBasicResponseDto](palette),
	}, err
}

func (c *PaletteService) Create(dto *request.PaletteDto, body *request.PaletteCreateDto) (*common.UuidResponseDto, error) {
	palette, err := c.Repository.Create(nil, dto, body)
	return common.NewResponse[common.UuidResponseDto](palette), err
}

func (c *PaletteService) Update(dto *request.PaletteDto, body *request.PaletteCreateDto) (*common.UuidResponseDto, error) {
	palette, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.Repository.Update(nil, dto, body, palette)
	return common.NewResponse[common.UuidResponseDto](entity), err
}

func (c *PaletteService) Delete(dto *request.PaletteDto) (interface{}, error) {
	palette, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	return nil, c.Repository.Delete(nil, dto, palette)
}

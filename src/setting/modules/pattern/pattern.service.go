package pattern

import (
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/setting/modules/pattern/dto/request"
	"grape/src/setting/modules/pattern/dto/response"
	repo "grape/src/setting/modules/pattern/repositories"
)

type PatternService struct {
	Repository *repo.PatternRepositoryT
}

func NewPatternService(s *service.CommonService) *PatternService {
	return &PatternService{
		Repository: repo.NewPatternRepository(s.DB),
	}
}

func (c *PatternService) FindOne(dto *request.PatternDto) (*response.PatternBasicResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto)

	return common.NewResponse[response.PatternBasicResponseDto](project), err
}

func (c *PatternService) FindAll(dto *request.PatternDto) (*common.PageResponseDto[[]response.PatternBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(dto)

	return &common.PageResponseDto[[]response.PatternBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.PatternBasicResponseDto](projects),
	}, err
}

func (c *PatternService) Create(dto *request.PatternDto, body *request.PatternCreateDto) (*common.UuidResponseDto, error) {
	project, err := c.Repository.Create(nil, dto, body)
	return common.NewResponse[common.UuidResponseDto](project), err
}

func (c *PatternService) Update(dto *request.PatternDto, body *request.PatternUpdateDto) (*common.UuidResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.Repository.Update(nil, dto, body, project)
	return common.NewResponse[common.UuidResponseDto](entity), err
}

func (c *PatternService) Delete(dto *request.PatternDto) (interface{}, error) {
	project, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	return nil, c.Repository.Delete(nil, dto, project)
}

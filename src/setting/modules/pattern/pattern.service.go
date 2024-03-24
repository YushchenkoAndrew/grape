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

func (c *PatternService) FindOne(dto *request.PatternDto) (*response.PatternAdvancedResponseDto, error) {
	pattern, err := c.Repository.ValidateEntityExistence(dto)

	return common.NewResponse[response.PatternAdvancedResponseDto](pattern), err
}

func (c *PatternService) FindAll(dto *request.PatternDto) (*common.PageResponseDto[[]response.PatternBasicResponseDto], error) {
	total, patterns, err := c.Repository.GetAllPage(dto)

	return &common.PageResponseDto[[]response.PatternBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.PatternBasicResponseDto](patterns),
	}, err
}

func (c *PatternService) Create(dto *request.PatternDto, body *request.PatternCreateDto) (*common.UuidResponseDto, error) {
	pattern, err := c.Repository.Create(nil, dto, body)
	return common.NewResponse[common.UuidResponseDto](pattern), err
}

func (c *PatternService) Update(dto *request.PatternDto, body *request.PatternUpdateDto) (*common.UuidResponseDto, error) {
	pattern, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.Repository.Update(nil, dto, body, pattern)
	return common.NewResponse[common.UuidResponseDto](entity), err
}

func (c *PatternService) Delete(dto *request.PatternDto) (interface{}, error) {
	pattern, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	return nil, c.Repository.Delete(nil, dto, pattern)
}

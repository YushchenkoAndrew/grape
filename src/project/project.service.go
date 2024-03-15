package project

import (
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/project/dto/request"
	"grape/src/project/dto/response"
	repo "grape/src/project/repositories"
	"grape/src/project/types"
)

type ProjectService struct {
	Repository *repo.ProjectRepositoryT

	// Link         i.Default[m.Link, m.LinkQueryDto]
	// File         i.Default[m.File, m.FileQueryDto]
	// Project      i.Default[m.Project, m.ProjectQueryDto]
	// Cron         *CronService
	// Subscription i.Default[m.Subscription, m.SubscribeQueryDto]
	// Metrics      i.Default[m.Metrics, m.MetricsQueryDto]
}

func NewProjectService(s *service.CommonService) *ProjectService {
	return &ProjectService{
		Repository: repo.NewProjectRepository(s.DB),
		// Link:         NewLinkService(db, client),
		// File:         NewFileService(db, client),
		// Project:      NewProjectService(db, client),
		// Cron:         NewCronService(),
		// Subscription: NewSubscriptionService(db, client),
		// Metrics:      pods.NewMetricsService(db, client),
	}
}

func (c *ProjectService) FindOne(dto *request.ProjectDto) (*response.ProjectBasicResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	}), repo.Attachments)

	return common.NewResponse[response.ProjectBasicResponseDto](project), err
}

func (c *ProjectService) FindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.ProjectBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	}))

	return &common.PageResponseDto[[]response.ProjectBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.ProjectBasicResponseDto](projects),
	}, err
}

func (c *ProjectService) AdminFindOne(dto *request.ProjectDto) (*response.ProjectAdvancedResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto, repo.Owner, repo.Attachments)

	return common.NewResponse[response.ProjectAdvancedResponseDto](project), err
}

func (c *ProjectService) AdminFindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.ProjectAdvancedResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(dto, repo.Owner)

	return &common.PageResponseDto[[]response.ProjectAdvancedResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.ProjectAdvancedResponseDto](projects),
	}, err
}

func (c *ProjectService) Create(dto *request.ProjectDto, body *request.ProjectCreateDto) (*common.UuidResponseDto, error) {
	entity, err := c.Repository.Create(nil, dto, body)

	return common.NewResponse[common.UuidResponseDto](entity), err
}

func (c *ProjectService) Update(dto *request.ProjectDto, body *request.ProjectUpdateDto) (*common.UuidResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.Repository.Update(nil, dto, body, project)
	return common.NewResponse[common.UuidResponseDto](entity), err
}

func (c *ProjectService) Delete(dto *request.ProjectDto, body *request.ProjectUpdateDto) (interface{}, error) {
	project, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	return nil, c.Repository.Delete(nil, dto, project)
}

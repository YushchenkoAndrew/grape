package project

import (
	"fmt"
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/project/dto/request"
	"grape/src/project/dto/response"
	e "grape/src/project/entities"
	repo "grape/src/project/repositories"
	"grape/src/project/types"

	"github.com/google/uuid"
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

	// (&repo.ProjectRepository{}).GetAll()
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
	project, err := c.ValidateProjectExistence(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	}))

	if err != nil {
		return nil, err
	}

	return common.NewResponse[response.ProjectBasicResponseDto](project), nil
}

func (c *ProjectService) FindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.ProjectBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	}))

	if err != nil {
		return nil, err
	}

	return &common.PageResponseDto[[]response.ProjectBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.ProjectBasicResponseDto](projects),
	}, nil
}

func (c *ProjectService) AdminFindOne(dto *request.ProjectDto) (*response.ProjectAdvancedResponseDto, error) {
	project, err := c.ValidateProjectExistence(dto, repo.Owner, repo.Attachments)
	if err != nil {
		return nil, err
	}

	return common.NewResponse[response.ProjectAdvancedResponseDto](project), nil
}

func (c *ProjectService) AdminFindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.ProjectAdvancedResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(dto, repo.Owner)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v", projects[0].Owner)

	return &common.PageResponseDto[[]response.ProjectAdvancedResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.ProjectAdvancedResponseDto](projects),
	}, nil
}

func (c *ProjectService) Create(dto *request.ProjectDto, body *request.ProjectCreateDto) (*common.UuidResponseDto, error) {
	entity, err := c.Repository.Create(dto, body)
	if err != nil {
		return nil, err
	}

	return common.NewResponse[common.UuidResponseDto](entity), nil
}

func (c *ProjectService) Update(dto *request.ProjectDto, body *request.ProjectUpdateDto) (*common.UuidResponseDto, error) {
	project, err := c.ValidateProjectExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.Repository.Update(dto, body, project)
	if err != nil {
		return nil, err
	}

	return common.NewResponse[common.UuidResponseDto](entity), nil
}

func (c *ProjectService) Delete(dto *request.ProjectDto, body *request.ProjectUpdateDto) (interface{}, error) {
	project, err := c.ValidateProjectExistence(dto)
	if err != nil {
		return nil, err
	}

	if err := c.Repository.Delete(dto, project); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *ProjectService) ValidateProjectExistence(dto *request.ProjectDto, relations ...repo.ProjectRelation) (*e.ProjectEntity, error) {
	if uuid.Validate(dto.ProjectIds[0]) != nil {
		return nil, fmt.Errorf("project id is invalid")
	}

	if project, _ := c.Repository.GetOne(dto, relations...); project != nil {
		return project, nil
	}

	return nil, fmt.Errorf("project not found")
}

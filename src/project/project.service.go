package project

import (
	"fmt"
	att "grape/src/attachment"
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/project/dto/request"
	"grape/src/project/dto/response"
	"grape/src/project/entities"
	repo "grape/src/project/repositories"
	"grape/src/project/types"

	"gorm.io/gorm"
)

type ProjectService struct {
	Repository *repo.ProjectRepositoryT

	AttachmentService *att.AttachmentService

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

		AttachmentService: att.NewAttachmentService(s),
		// Link:         NewLinkService(db, client),
		// File:         NewFileService(db, client),
		// Project:      NewProjectService(db, client),
		// Cron:         NewCronService(),
		// Subscription: NewSubscriptionService(db, client),
		// Metrics:      pods.NewMetricsService(db, client),
	}
}

func (c *ProjectService) FindOne(dto *request.ProjectDto) (*response.ProjectDetailedResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	}), repo.Attachments, repo.ColorPalette, repo.SvgPattern)

	return common.NewResponse[response.ProjectDetailedResponseDto](project), err
}

func (c *ProjectService) FindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.ProjectBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	}), repo.Attachments, repo.ColorPalette, repo.SvgPattern)

	return &common.PageResponseDto[[]response.ProjectBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.ProjectBasicResponseDto](projects),
	}, err
}

func (c *ProjectService) AdminFindOne(dto *request.ProjectDto) (*response.AdminProjectDetailedResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto, repo.Owner, repo.Attachments, repo.ColorPalette, repo.SvgPattern, repo.Statistic)

	return common.NewResponse[response.AdminProjectDetailedResponseDto](project), err
}

func (c *ProjectService) AdminFindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.AdminProjectBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(dto, repo.Owner, repo.Attachments, repo.ColorPalette, repo.SvgPattern, repo.Statistic)

	return &common.PageResponseDto[[]response.AdminProjectBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.AdminProjectBasicResponseDto](projects),
	}, err
}

func (c *ProjectService) Create(dto *request.ProjectDto, body *request.ProjectCreateDto) (*common.UuidResponseDto, error) {
	var project *entities.ProjectEntity
	err := c.Repository.Transaction(func(tx *gorm.DB) error {
		var err error
		project, err = c.Repository.Create(tx, dto, body)
		if err != nil {
			return err
		}

		attachments := c.AttachmentService.InitProjectFromTemplate(project)
		if err := tx.Model(project).Association("Attachments").Append(&attachments); err != nil {
			return err
		}

		return nil
	})

	return common.NewResponse[common.UuidResponseDto](project), err
}

func (c *ProjectService) Update(dto *request.ProjectDto, body *request.ProjectUpdateDto) (*common.UuidResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.Repository.Update(nil, dto, body, project)
	return common.NewResponse[common.UuidResponseDto](entity), err
}

func (c *ProjectService) Delete(dto *request.ProjectDto) (interface{}, error) {
	project, err := c.Repository.ValidateEntityExistence(dto, repo.Attachments)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		if _, err := c.AttachmentService.VoidService.Delete(project.GetPath()); err != nil {
			fmt.Println(err.Error())
			return err
		}

		for _, attachment := range project.Attachments {
			if err := c.AttachmentService.Repository.Delete(tx, nil, &attachment); err != nil {
				return nil
			}
		}

		return c.Repository.Delete(nil, dto, project)
	})

	return nil, err
}

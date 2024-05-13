package project

import (
	"grape/src/attachment"
	att_req "grape/src/attachment/dto/request"
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	t "grape/src/common/types"
	ln_entity "grape/src/link/entities"
	ln_repo "grape/src/link/repositories"
	"grape/src/project/dto/request"
	"grape/src/project/dto/response"
	"grape/src/project/entities"
	repo "grape/src/project/repositories"
	"grape/src/project/types"
	st_req "grape/src/statistic/dto/request"
	st_repo "grape/src/statistic/repositories"

	"gorm.io/gorm"
)

type ProjectService struct {
	Repository          *repo.ProjectRepositoryT
	StatisticRepository *st_repo.StatisticRepositoryT
	LinkRepository      *ln_repo.LinkRepositoryT

	AttachmentService *attachment.AttachmentService

	// Cron         *CronService
	// Subscription i.Default[m.Subscription, m.SubscribeQueryDto]
	// Metrics      i.Default[m.Metrics, m.MetricsQueryDto]
}

func NewProjectService(s *service.CommonService) *ProjectService {
	return &ProjectService{
		Repository:          repo.NewProjectRepository(s.DB),
		LinkRepository:      ln_repo.NewLinkRepository(s.DB),
		StatisticRepository: st_repo.NewStatisticRepository(s.DB),

		AttachmentService: attachment.NewAttachmentService(s),
		// Cron:         NewCronService(),
		// Subscription: NewSubscriptionService(db, client),
		// Metrics:      pods.NewMetricsService(db, client),
	}
}

func (c *ProjectService) FindOne(dto *request.ProjectDto) (*response.ProjectDetailedResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{t.Active.String()},
	}), repo.Attachments, repo.Links)

	return common.NewResponse[response.ProjectDetailedResponseDto](project), err
}

func (c *ProjectService) FindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.ProjectBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{t.Active.String()},
	}), repo.Attachments)

	return &common.PageResponseDto[[]response.ProjectBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.ProjectBasicResponseDto](projects),
	}, err
}

func (c *ProjectService) AdminFindOne(dto *request.ProjectDto) (*response.AdminProjectDetailedResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto, repo.Owner, repo.Attachments, repo.Statistic, repo.Links)

	return common.NewResponse[response.AdminProjectDetailedResponseDto](project), err
}

func (c *ProjectService) AdminFindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.AdminProjectBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(dto, repo.Owner, repo.Attachments, repo.Links, repo.Statistic)

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

		if body.Type == types.Link.String() {
			entity := ln_entity.NewLinkEntity()
			entity.Name, entity.Link = "redirect", body.Link
			entity.Create()

			if err = tx.Model(project).Association("Links").Append(entity); err != nil {
				return err
			}
		}

		attachments := c.AttachmentService.InitProjectFromTemplate(project, body.README)
		if err = tx.Model(project).Association("Attachments").Append(&attachments); err != nil {
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
	project, err := c.Repository.ValidateEntityExistence(dto, repo.Attachments, repo.Links)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		for _, e := range project.Attachments {
			dto := att_req.NewAttachmentDto(dto.CurrentUser, &att_req.AttachmentDto{AttachmentIds: []string{e.UUID}})
			if _, err := c.AttachmentService.Delete(dto); err != nil {
				return nil
			}
		}

		if err := c.LinkRepository.DeleteAll(tx, nil, project.GetLinks()); err != nil {
			return nil
		}

		return c.Repository.Delete(tx, dto, project)
	})

	return nil, err
}

func (c *ProjectService) UpdateProjectStatistics(dto *request.ProjectDto, body *st_req.StatisticUpdateDto) (interface{}, error) {
	project, err := c.Repository.ValidateEntityExistence(dto, repo.Statistic)
	if err != nil {
		return nil, err
	}

	_, err = c.StatisticRepository.Update(nil, nil, body, project.Statistic)
	return nil, err
}

func (c *ProjectService) UpdateOrder(dto *request.ProjectDto, body *req.OrderUpdateDto) (*common.UuidResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil || project.Order == body.Position {
		return common.NewResponse[common.UuidResponseDto](project), err
	}

	err = c.Repository.Reorder(nil, project, body.Position)
	return common.NewResponse[common.UuidResponseDto](project), err
}

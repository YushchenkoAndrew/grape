package project

import (
	att "grape/src/attachment"
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	ln_entity "grape/src/link/entities"
	ln_repo "grape/src/link/repositories"
	"grape/src/project/dto/request"
	"grape/src/project/dto/response"
	"grape/src/project/entities"
	repo "grape/src/project/repositories"
	"grape/src/project/types"
	statistic "grape/src/statistic/dto/request"
	st_repo "grape/src/statistic/repositories"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ProjectService struct {
	Repository          *repo.ProjectRepositoryT
	StatisticRepository *st_repo.StatisticRepositoryT

	LinkRepository    *ln_repo.LinkRepositoryT
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
		Repository:          repo.NewProjectRepository(s.DB),
		LinkRepository:      ln_repo.NewLinkRepository(s.DB),
		StatisticRepository: st_repo.NewStatisticRepository(s.DB),

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
	}), repo.Attachments, repo.Palette, repo.Pattern, repo.Links)

	return common.NewResponse[response.ProjectDetailedResponseDto](project), err
}

func (c *ProjectService) FindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.ProjectBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	}), repo.Attachments, repo.Palette, repo.Pattern)

	return &common.PageResponseDto[[]response.ProjectBasicResponseDto]{
		Page:    dto.Page,
		PerPage: dto.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.ProjectBasicResponseDto](projects),
	}, err
}

func (c *ProjectService) AdminFindOne(dto *request.ProjectDto) (*response.AdminProjectDetailedResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto, repo.Owner, repo.Attachments, repo.Palette, repo.Pattern, repo.Statistic, repo.Links)

	return common.NewResponse[response.AdminProjectDetailedResponseDto](project), err
}

func (c *ProjectService) AdminFindAll(dto *request.ProjectDto) (*common.PageResponseDto[[]response.AdminProjectBasicResponseDto], error) {
	total, projects, err := c.Repository.GetAllPage(dto, repo.Owner, repo.Attachments, repo.Palette, repo.Pattern, repo.Statistic)

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
		var projects []*entities.ProjectEntity
		res := tx.Table(c.Repository.TableName()).
			Where(`projects.organization_id = ?`, dto.CurrentUser.Organization.ID).
			Where(`projects.order > ?`, project.Order).
			Find(&projects)

		if res.Error != nil {
			return res.Error
		}

		if len(projects) != 0 {
			lo.ForEach(projects, func(e *entities.ProjectEntity, _ int) { e.Order -= 1 })
			if res := tx.Table(c.Repository.TableName()).Save(projects); res.Error != nil {
				return res.Error
			}
		}

		if _, err := c.AttachmentService.VoidService.Delete(project.GetPath()); err != nil {
			return err
		}

		for _, attachment := range project.Attachments {
			if err := c.AttachmentService.Repository.Delete(tx, nil, &attachment); err != nil {
				return nil
			}
		}

		for _, link := range project.Links {
			if err := c.LinkRepository.Delete(tx, nil, &link); err != nil {
				return nil
			}
		}

		return c.Repository.Delete(nil, dto, project)
	})

	return nil, err
}

func (c *ProjectService) UpdateProjectStatistics(dto *request.ProjectDto, body *statistic.StatisticUpdateDto) (interface{}, error) {
	project, err := c.Repository.ValidateEntityExistence(dto, repo.Statistic)
	if err != nil {
		return nil, err
	}

	_, err = c.StatisticRepository.Update(nil, nil, body, project.Statistic)
	return nil, err
}

func (c *ProjectService) PutOrder(dto *request.ProjectDto, body *request.ProjectOrderUpdateDto) (*common.UuidResponseDto, error) {
	project, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil || project.Order == body.Position {
		return common.NewResponse[common.UuidResponseDto](project), err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		var projects []*entities.ProjectEntity
		db := tx.Table(c.Repository.TableName()).Where(`projects.organization_id = ?`, dto.CurrentUser.Organization.ID)

		if project.Order < body.Position {
			db = db.Where(`projects.order > ?`, project.Order).Where(`projects.order <= ?`, body.Position)
		} else {
			db = db.Where(`projects.order < ?`, project.Order).Where(`projects.order >= ?`, body.Position)
		}

		if res := db.Find(&projects); res.Error != nil || len(projects) == 0 {
			return res.Error
		}

		for _, e := range projects {
			if project.Order < body.Position {
				e.Order -= 1
			} else {
				e.Order += 1
			}
		}

		project.Order = body.Position
		projects = append(projects, project)

		res := tx.Table(c.Repository.TableName()).Save(projects)
		return res.Error
	})

	return common.NewResponse[common.UuidResponseDto](project), err
}

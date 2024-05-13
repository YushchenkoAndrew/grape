package stage

import (
	"grape/src/attachment"
	att_req "grape/src/attachment/dto/request"
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/common/types"
	ctx_entity "grape/src/context/entities"
	ctx_repo "grape/src/context/repositories"
	ln_entity "grape/src/link/entities"
	ln_repo "grape/src/link/repositories"
	"grape/src/stage/dto/request"
	"grape/src/stage/dto/response"
	repo "grape/src/stage/repositories"

	"gorm.io/gorm"
)

type StageService struct {
	Repository        *repo.StageRepositoryT
	TaskRepository    *repo.TaskRepositoryT
	LinkRepository    *ln_repo.LinkRepositoryT
	ContextRepository *ctx_repo.ContextRepositoryT

	AttachmentService *attachment.AttachmentService
}

func NewStageService(s *service.CommonService) *StageService {
	return &StageService{
		Repository:        repo.NewStageRepository(s.DB),
		TaskRepository:    repo.NewTaskRepository(s.DB),
		LinkRepository:    ln_repo.NewLinkRepository(s.DB),
		ContextRepository: ctx_repo.NewContextRepository(s.DB),

		AttachmentService: attachment.NewAttachmentService(s),
	}
}

func (c *StageService) FindAll(dto *request.StageDto) ([]response.StageBasicResponseDto, error) {
	stages, err := c.Repository.GetAll(req.NewRequest(dto, &request.StageDto{
		Statuses: []string{types.Active.String()},
	}), repo.Tasks)

	return common.NewResponseMany[response.StageBasicResponseDto](stages), err
}

func (c *StageService) AdminFindAll(dto *request.StageDto) ([]response.AdminStageBasicResponseDto, error) {
	stages, err := c.Repository.GetAll(dto, repo.Tasks)

	return common.NewResponseMany[response.AdminStageBasicResponseDto](stages), err
}

func (c *StageService) Create(dto *request.StageDto, body *request.StageCreateDto) (*common.UuidResponseDto, error) {
	stage, err := c.Repository.Create(nil, dto, body)
	return common.NewResponse[common.UuidResponseDto](stage), err
}

func (c *StageService) Update(dto *request.StageDto, body *request.StageUpdateDto) (*common.UuidResponseDto, error) {
	stage, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.Repository.Update(nil, dto, body, stage)
	return common.NewResponse[common.UuidResponseDto](entity), err
}

func (c *StageService) Delete(dto *request.StageDto) (interface{}, error) {
	stage, err := c.Repository.ValidateEntityExistence(dto, repo.Tasks)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		var links []*ln_entity.LinkEntity
		var contexts []*ctx_entity.ContextEntity

		for _, task := range stage.Tasks {
			links = append(links, task.GetLinks()...)
			contexts = append(contexts, task.GetContexts()...)

			for _, e := range task.Attachments {
				dto := att_req.NewAttachmentDto(dto.CurrentUser, &att_req.AttachmentDto{AttachmentIds: []string{e.UUID}})
				if _, err := c.AttachmentService.Delete(dto); err != nil {
					return nil
				}
			}
		}

		if err := c.LinkRepository.DeleteAll(tx, nil, links); err != nil {
			return nil
		}

		if err := c.ContextRepository.DeleteAll(tx, nil, contexts); err != nil {
			return nil
		}

		if err := c.TaskRepository.DeleteAll(tx, nil, stage.GetTasks()); err != nil {
			return nil
		}

		return c.Repository.Delete(tx, dto, stage)
	})

	return nil, err
}

func (c *StageService) CreateTask(dto *request.TaskDto, body *request.TaskCreateDto) (*common.UuidResponseDto, error) {
	stage, err := c.Repository.ValidateEntityExistence(request.NewStageDto(dto.CurrentUser, &request.StageDto{StageIds: dto.StageIds}))
	if err != nil {
		return nil, err
	}

	body.Stage = stage
	task, err := c.TaskRepository.Create(nil, dto, body)
	return common.NewResponse[common.UuidResponseDto](task), err
}

func (c *StageService) UpdateTask(dto *request.TaskDto, body *request.TaskUpdateDto) (*common.UuidResponseDto, error) {
	task, err := c.TaskRepository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.TaskRepository.Update(nil, dto, body, task)
	return common.NewResponse[common.UuidResponseDto](entity), err
}

func (c *StageService) DeleteTask(dto *request.TaskDto) (interface{}, error) {
	task, err := c.TaskRepository.ValidateEntityExistence(dto, repo.Attachments, repo.Links, repo.Contexts)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		for _, e := range task.Attachments {
			dto := att_req.NewAttachmentDto(dto.CurrentUser, &att_req.AttachmentDto{AttachmentIds: []string{e.UUID}})
			if _, err := c.AttachmentService.Delete(dto); err != nil {
				return nil
			}
		}

		if err := c.LinkRepository.DeleteAll(tx, nil, task.GetLinks()); err != nil {
			return nil
		}

		if err := c.ContextRepository.DeleteAll(tx, nil, task.GetContexts()); err != nil {
			return nil
		}

		return c.TaskRepository.Delete(tx, dto, task)
	})

	return nil, err
}

func (c *StageService) UpdateOrder(dto *request.StageDto, body *req.OrderUpdateDto) (*common.UuidResponseDto, error) {
	stage, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil || stage.Order == body.Position {
		return common.NewResponse[common.UuidResponseDto](stage), err
	}

	err = c.Repository.Reorder(nil, stage, body.Position)
	return common.NewResponse[common.UuidResponseDto](stage), err
}

func (c *StageService) UpdateTaskOrder(dto *request.TaskDto, body *req.OrderUpdateDto) (*common.UuidResponseDto, error) {
	task, err := c.TaskRepository.ValidateEntityExistence(dto)
	if err != nil || task.Order == body.Position {
		return common.NewResponse[common.UuidResponseDto](task), err
	}

	err = c.TaskRepository.Reorder(nil, task, body.Position)
	return common.NewResponse[common.UuidResponseDto](task), err
}

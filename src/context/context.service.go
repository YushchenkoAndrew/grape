package context

import (
	"fmt"
	req "grape/src/common/dto/request"
	res "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/context/dto/request"
	"grape/src/context/dto/response"
	"grape/src/context/entities"
	"grape/src/context/repositories"

	st_req "grape/src/stage/dto/request"
	st_repo "grape/src/stage/repositories"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ContextService struct {
	Repository             *repositories.ContextRepositoryT
	ContextFieldRepository *repositories.ContextFieldRepositoryT
	TaskRepository         *st_repo.TaskRepositoryT
}

func NewContextService(s *service.CommonService) *ContextService {
	return &ContextService{
		Repository:             repositories.NewContextRepository(s.DB),
		ContextFieldRepository: repositories.NewContextFieldRepository(s.DB),
		TaskRepository:         st_repo.NewTaskRepository(s.DB),
	}
}

func (c *ContextService) AdminFindOne(dto *request.ContextDto) (*response.ContextAdvancedResponseDto, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto, repositories.ContextFields)
	return res.NewResponse[response.ContextAdvancedResponseDto](entity), err
}

func (c *ContextService) Create(dto *request.ContextDto, body *request.ContextCreateDto) (*res.UuidResponseDto, error) {
	entity := entities.NewContextEntity()
	copier.Copy(&entity, body)
	entity.Create()

	contextable, err := func() (entities.ContextableT, error) {
		switch body.ContextableType {
		case c.TaskRepository.TableName():
			return c.TaskRepository.ValidateEntityExistence(st_req.NewTaskDto(dto.CurrentUser, &st_req.TaskDto{TaskIds: []string{body.ContextableID}}))
		}

		return nil, fmt.Errorf("contextable_type '%s' is not supported", body.ContextableType)
	}()

	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		entity.ContextableID = contextable.GetID()
		entity.ContextableType = body.ContextableType
		_, err := c.Repository.Create(tx, dto, entity)
		return err
	})

	return res.NewResponse[res.UuidResponseDto](entity), err
}

func (c *ContextService) Update(dto *request.ContextDto, body *request.ContextUpdateDto) (*res.UuidResponseDto, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	_, err = c.Repository.Update(nil, dto, body, entity)
	return res.NewResponse[res.UuidResponseDto](entity), err
}

func (c *ContextService) Delete(dto *request.ContextDto) (interface{}, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		if err := c.ContextFieldRepository.DeleteAll(tx, nil, entity.GetContextFields()); err != nil {
			return nil
		}

		return c.Repository.Delete(tx, dto, entity)
	})

	return nil, err
}

func (c *ContextService) CreateField(dto *request.ContextFieldDto, body *request.ContextFieldCreateDto) (*res.UuidResponseDto, error) {
	context, err := c.Repository.ValidateEntityExistence(request.NewContextDto(dto.CurrentUser, &request.ContextDto{ContextIds: dto.ContextIds}))
	if err != nil {
		return nil, err
	}

	body.Context = context
	field, err := c.ContextFieldRepository.Create(nil, dto, body)
	return res.NewResponse[res.UuidResponseDto](field), err
}

func (c *ContextService) UpdateField(dto *request.ContextFieldDto, body *request.ContextFieldUpdateDto) (*res.UuidResponseDto, error) {
	field, err := c.ContextFieldRepository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	entity, err := c.ContextFieldRepository.Update(nil, dto, body, field)
	return res.NewResponse[res.UuidResponseDto](entity), err
}

func (c *ContextService) DeleteField(dto *request.ContextFieldDto) (interface{}, error) {
	field, err := c.ContextFieldRepository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	err = c.ContextFieldRepository.Delete(nil, dto, field)
	return nil, err
}

func (c *ContextService) UpdateOrder(dto *request.ContextDto, body *req.OrderUpdateDto) (*res.UuidResponseDto, error) {
	link, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil || link.Order == body.Position {
		return res.NewResponse[res.UuidResponseDto](link), err
	}

	err = c.Repository.Reorder(nil, link, body.Position)
	return res.NewResponse[res.UuidResponseDto](link), err
}

func (c *ContextService) UpdateFieldOrder(dto *request.ContextFieldDto, body *req.OrderUpdateDto) (*res.UuidResponseDto, error) {
	field, err := c.ContextFieldRepository.ValidateEntityExistence(dto)
	if err != nil || field.Order == body.Position {
		return res.NewResponse[res.UuidResponseDto](field), err
	}

	err = c.ContextFieldRepository.Reorder(nil, field, body.Position)
	return res.NewResponse[res.UuidResponseDto](field), err
}

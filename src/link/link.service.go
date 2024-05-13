package link

import (
	"fmt"
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/link/dto/request"
	"grape/src/link/dto/response"
	"grape/src/link/entities"
	repo "grape/src/link/repositories"

	pr_req "grape/src/project/dto/request"
	pr_repo "grape/src/project/repositories"

	st_req "grape/src/stage/dto/request"
	st_repo "grape/src/stage/repositories"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type LinkService struct {
	Repository        *repo.LinkRepositoryT
	TaskRepository    *st_repo.TaskRepositoryT
	ProjectRepository *pr_repo.ProjectRepositoryT
}

func NewLinkService(s *service.CommonService) *LinkService {
	return &LinkService{
		Repository:        repo.NewLinkRepository(s.DB),
		TaskRepository:    st_repo.NewTaskRepository(s.DB),
		ProjectRepository: pr_repo.NewProjectRepository(s.DB),
	}
}

func (c *LinkService) AdminFindOne(dto *request.LinkDto) (*response.LinkAdvancedResponseDto, error) {
	link, err := c.Repository.ValidateEntityExistence(dto)
	return common.NewResponse[response.LinkAdvancedResponseDto](link), err
}

func (c *LinkService) Create(dto *request.LinkDto, body *request.LinkCreateDto) (*response.LinkBasicResponseDto, error) {
	entity := entities.NewLinkEntity()
	copier.Copy(&entity, body)
	entity.Create()

	linkable, err := func() (entities.LinkableT, error) {
		switch body.LinkableType {
		case c.ProjectRepository.TableName():
			return c.ProjectRepository.ValidateEntityExistence(pr_req.NewProjectDto(dto.CurrentUser, &pr_req.ProjectDto{ProjectIds: []string{body.LinkableID}}))

		case c.TaskRepository.TableName():
			return c.TaskRepository.ValidateEntityExistence(st_req.NewTaskDto(dto.CurrentUser, &st_req.TaskDto{TaskIds: []string{body.LinkableID}}))
		}

		return nil, fmt.Errorf("linkable_type '%s' is not supported", body.LinkableType)
	}()

	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		entity.LinkableID = linkable.GetID()
		entity.LinkableType = body.LinkableType
		_, err := c.Repository.Create(tx, dto, entity)
		return err
	})

	return common.NewResponse[response.LinkBasicResponseDto](entity), err
}

func (c *LinkService) Update(dto *request.LinkDto, body *request.LinkUpdateDto) (*response.LinkBasicResponseDto, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	_, err = c.Repository.Update(nil, dto, body, entity)
	return common.NewResponse[response.LinkBasicResponseDto](entity), err
}

func (c *LinkService) Delete(dto *request.LinkDto) (interface{}, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Delete(nil, dto, entity)
	return nil, err
}

func (c *LinkService) UpdateOrder(dto *request.LinkDto, body *req.OrderUpdateDto) (*response.LinkBasicResponseDto, error) {
	link, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil || link.Order == body.Position {
		return common.NewResponse[response.LinkBasicResponseDto](link), err
	}

	err = c.Repository.Reorder(nil, link, body.Position)
	return common.NewResponse[response.LinkBasicResponseDto](link), err
}

package tag

import (
	"fmt"
	"grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/tag/dto/request"
	"grape/src/tag/entities"
	repo "grape/src/tag/repositories"

	pr_req "grape/src/project/dto/request"
	pr_repo "grape/src/project/repositories"

	st_req "grape/src/stage/dto/request"
	st_repo "grape/src/stage/repositories"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type TagService struct {
	Repository        *repo.TagRepositoryT
	TaskRepository    *st_repo.TaskRepositoryT
	ProjectRepository *pr_repo.ProjectRepositoryT
}

func NewTagService(s *service.CommonService) *TagService {
	return &TagService{
		Repository:        repo.NewTagRepository(s.DB),
		TaskRepository:    st_repo.NewTaskRepository(s.DB),
		ProjectRepository: pr_repo.NewProjectRepository(s.DB),
	}
}

func (c *TagService) AdminFindOne(dto *request.TagDto) (*response.UuidResponseDto, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	return response.NewResponse[response.UuidResponseDto](entity), err
}

func (c *TagService) Create(dto *request.TagDto, body *request.TagCreateDto) (*response.UuidResponseDto, error) {
	entity := entities.NewTagEntity()
	copier.Copy(&entity, body)
	entity.Create()

	linkable, err := func() (entities.TaggableT, error) {
		switch body.TaggableType {
		case c.ProjectRepository.TableName():
			return c.ProjectRepository.ValidateEntityExistence(pr_req.NewProjectDto(dto.CurrentUser, &pr_req.ProjectDto{ProjectIds: []string{body.TaggableID}}))

		case c.TaskRepository.TableName():
			return c.TaskRepository.ValidateEntityExistence(st_req.NewTaskDto(dto.CurrentUser, &st_req.TaskDto{TaskIds: []string{body.TaggableID}}))
		}

		return nil, fmt.Errorf("taggable_type '%s' is not supported", body.TaggableType)
	}()

	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		entity.TaggableID = linkable.GetID()
		entity.TaggableType = body.TaggableType
		_, err := c.Repository.Create(tx, dto, entity)
		return err
	})

	return response.NewResponse[response.UuidResponseDto](entity), err
}

func (c *TagService) Update(dto *request.TagDto, body *request.TagUpdateDto) (*response.UuidResponseDto, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	_, err = c.Repository.Update(nil, dto, body, entity)
	return response.NewResponse[response.UuidResponseDto](entity), err
}

func (c *TagService) Delete(dto *request.TagDto) (interface{}, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Delete(nil, dto, entity)
	return nil, err
}

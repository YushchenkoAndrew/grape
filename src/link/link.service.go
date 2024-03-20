package link

import (
	"fmt"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/link/dto/request"
	"grape/src/link/dto/response"
	"grape/src/link/entities"
	repo "grape/src/link/repositories"
	project "grape/src/project/dto/request"

	pr "grape/src/project/repositories"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type LinkService struct {
	Repository        *repo.LinkRepositoryT
	ProjectRepository *pr.ProjectRepositoryT
}

func NewLinkService(s *service.CommonService) *LinkService {
	return &LinkService{
		Repository:        repo.NewLinkRepository(s.DB),
		ProjectRepository: pr.NewProjectRepository(s.DB),
	}
}

func (c *LinkService) Create(dto *request.LinkDto, body *request.LinkCreateDto) (*response.LinkBasicResponseDto, error) {
	entity := entities.NewLinkEntity()
	copier.Copy(&entity, body)
	entity.Create()

	switch body.LinkableType {
	case c.ProjectRepository.TableName():
		project, err := c.ProjectRepository.ValidateEntityExistence(project.NewProjectDto(dto.CurrentUser, &project.ProjectDto{ProjectIds: []string{body.LinkableID}}))
		if err != nil {
			return nil, err
		}

		err = c.ProjectRepository.Transaction(func(tx *gorm.DB) error {
			return tx.Model(project).Association("Links").Append(entity)
		})

		return common.NewResponse[response.LinkBasicResponseDto](entity), err
	}

	return nil, fmt.Errorf("linkable_type '%s' is not supported", body.LinkableType)
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

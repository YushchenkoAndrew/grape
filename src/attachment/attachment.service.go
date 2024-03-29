package attachment

import (
	"fmt"
	"grape/src/attachment/dto/request"
	"grape/src/attachment/dto/response"
	"grape/src/attachment/entities"
	"grape/src/attachment/repositories"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	pre "grape/src/project/entities"
	prt "grape/src/project/types"
	"grape/src/void"
	"mime/multipart"
	"path/filepath"

	project "grape/src/project/dto/request"
	pr "grape/src/project/repositories"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type AttachmentService struct {
	Repository        *repositories.AttachmentRepositoryT
	ProjectRepository *pr.ProjectRepositoryT

	VoidService *void.VoidService
}

func NewAttachmentService(s *service.CommonService) *AttachmentService {
	return &AttachmentService{
		Repository:        repositories.NewAttachmentRepository(s.DB),
		ProjectRepository: pr.NewProjectRepository(s.DB),

		VoidService: void.NewVoidService(s),
	}
}

func (c *AttachmentService) FindOne(dto *request.AttachmentDto) (string, []byte, error) {
	attachment, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return "", nil, err
	}

	return c.VoidService.Get(attachment.GetFile())
}

func (c *AttachmentService) Create(dto *request.AttachmentDto, body *request.AttachmentCreateDto, file *multipart.FileHeader) (*response.AttachmentAdvancedResponseDto, error) {
	entity := entities.NewAttachmentEntity()
	entity.Name, entity.Path, entity.Size, entity.Type = file.Filename, body.Path, file.Size, filepath.Ext(file.Filename)

	entity.Create()

	switch body.AttachableType {
	case c.ProjectRepository.TableName():
		project, err := c.ProjectRepository.ValidateEntityExistence(project.NewProjectDto(dto.CurrentUser, &project.ProjectDto{ProjectIds: []string{body.AttachableID}}))
		if err != nil {
			return nil, err
		}

		entity.Home = project.GetPath()
		err = c.ProjectRepository.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(project).Association("Attachments").Append(entity); err != nil {
				return err
			}

			f, err := file.Open()
			if err != nil {
				return err
			}

			defer f.Close()
			return c.VoidService.Save(entity.GetPath(), file.Filename, f)
		})

		return common.NewResponse[response.AttachmentAdvancedResponseDto](entity), err

	}

	return nil, fmt.Errorf("attachable_type '%s' is not supported", body.AttachableType)
}

func (c *AttachmentService) Update(dto *request.AttachmentDto, body *request.AttachmentUpdateDto, file *multipart.FileHeader) (*response.AttachmentAdvancedResponseDto, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		path := entity.GetFile()

		if file == nil {
			if _, err := c.Repository.Update(tx, dto, body, entity); err != nil {
				return err
			}

			return c.VoidService.Rename(path, entity.GetFile(), false)
		}

		e := entities.AttachmentEntity{}
		e.Name, e.Path, e.Size, e.Type = file.Filename, body.Path, file.Size, filepath.Ext(file.Filename)
		copier.CopyWithOption(&e, body, copier.Option{IgnoreEmpty: true})

		if _, err := c.Repository.Update(tx, dto, e, entity); err != nil {
			return err
		}

		f, err := file.Open()
		if err != nil {
			return err
		}

		defer f.Close()
		if err := c.VoidService.Save(entity.GetPath(), e.Name, f); err != nil {
			return err
		}

		_, err = c.VoidService.Delete(path)
		return err
	})

	return common.NewResponse[response.AttachmentAdvancedResponseDto](entity), err
}

func (c *AttachmentService) Delete(dto *request.AttachmentDto) (interface{}, error) {
	entity, err := c.Repository.ValidateEntityExistence(dto)
	if err != nil {
		return nil, err
	}

	err = c.Repository.Transaction(func(tx *gorm.DB) error {
		if err := c.Repository.Delete(tx, dto, entity); err != nil {
			return err
		}

		_, err = c.VoidService.Delete(entity.GetFile())
		return err

	})

	return nil, err
}

func (c *AttachmentService) InitProjectFromTemplate(project *pre.ProjectEntity) []entities.AttachmentEntity {
	var attachments []entities.AttachmentEntity

	entity := entities.NewAttachmentEntity()
	entity.Home = project.GetPath()
	entity.Create()

	switch project.Type {
	case prt.Html:
		entity.Name, entity.Path, entity.Size, entity.Type = "index.html", "/", 0, ".html"
		c.VoidService.Rename("/templates/html.template.html", entity.GetFile(), true)

	case prt.Markdown:
		entity.Name, entity.Path, entity.Size, entity.Type = "index.md", "/", 0, ".md"
		c.VoidService.Rename("/templates/markdown.template.md", entity.GetFile(), true)

	// TODO:
	// case K3s:
	default:
		entity = nil
	}

	if entity != nil {
		attachments = append(attachments, *entity)
	}

	return attachments
}

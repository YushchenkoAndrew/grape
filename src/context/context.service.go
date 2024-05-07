package context

import (
	"grape/src/common/service"
	"grape/src/context/repositories"
)

type ContextService struct {
	Repository *repositories.ContextRepositoryT
}

func NewContextService(s *service.CommonService) *ContextService {
	return &ContextService{
		Repository: repositories.NewContextRepository(s.DB),
	}
}

// func (c *ContextService) FindOne(dto *request.AttachmentDto) (string, []byte, error) {
// 	attachment, err := c.Repository.ValidateEntityExistence(dto)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	return c.VoidService.Get(attachment.GetFile())
// }

// func (c *ContextService) AdminFindOne(dto *request.AttachmentDto) (*response.AttachmentAdvancedResponseDto, error) {
// 	attachment, err := c.Repository.ValidateEntityExistence(dto)
// 	return common.NewResponse[response.AttachmentAdvancedResponseDto](attachment), err
// }

// func (c *ContextService) Create(dto *request.AttachmentDto, body *request.AttachmentCreateDto, file *multipart.FileHeader) (*response.AttachmentAdvancedResponseDto, error) {
// 	entity := entities.NewAttachmentEntity()
// 	entity.Name, entity.Path, entity.Size, entity.Type = file.Filename, body.Path, file.Size, filepath.Ext(file.Filename)

// 	entity.Create()

// 	switch body.AttachableType {
// 	case c.ProjectRepository.TableName():
// 		project, err := c.ProjectRepository.ValidateEntityExistence(project.NewProjectDto(dto.CurrentUser, &project.ProjectDto{ProjectIds: []string{body.AttachableID}}), pr.Attachments)
// 		if err != nil {
// 			return nil, err
// 		}

// 		entity.Home = project.GetPath()
// 		err = c.ProjectRepository.Transaction(func(tx *gorm.DB) error {
// 			entity.AttachableID = project.ID
// 			entity.AttachableType = c.ProjectRepository.TableName()
// 			if _, err := c.Repository.Create(tx, dto, entity); err != nil {
// 				return err
// 			}

// 			f, err := file.Open()
// 			if err != nil {
// 				return err
// 			}

// 			defer f.Close()
// 			return c.VoidService.Save(entity.GetPath(), file.Filename, f)
// 		})

// 		return common.NewResponse[response.AttachmentAdvancedResponseDto](entity), err

// 	}

// 	return nil, fmt.Errorf("attachable_type '%s' is not supported", body.AttachableType)
// }

// func (c *ContextService) Update(dto *request.AttachmentDto, body *request.AttachmentUpdateDto, file *multipart.FileHeader) (*response.AttachmentAdvancedResponseDto, error) {
// 	entity, err := c.Repository.ValidateEntityExistence(dto)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = c.Repository.Transaction(func(tx *gorm.DB) error {
// 		path := entity.GetFile()

// 		if file == nil {
// 			if _, err := c.Repository.Update(tx, dto, body, entity); err != nil {
// 				return err
// 			}

// 			return c.VoidService.Rename(path, entity.GetFile(), false)
// 		}

// 		e := entities.AttachmentEntity{}
// 		e.Name, e.Path, e.Home, e.Size, e.Type = file.Filename, body.Path, entity.Home, file.Size, filepath.Ext(file.Filename)
// 		copier.CopyWithOption(&e, body, copier.Option{IgnoreEmpty: true})
// 		updated, err := c.Repository.Update(tx, dto, e, entity)

// 		if err != nil {
// 			return err
// 		}

// 		if _, err := c.VoidService.Delete(path); err != nil {
// 			return err
// 		}

// 		f, err := file.Open()
// 		if err != nil {
// 			return err
// 		}

// 		defer f.Close()
// 		return c.VoidService.Save(updated.GetPath(), e.Name, f)
// 	})

// 	return common.NewResponse[response.AttachmentAdvancedResponseDto](entity), err
// }

// func (c *ContextService) Delete(dto *request.AttachmentDto) (interface{}, error) {
// 	entity, err := c.Repository.ValidateEntityExistence(dto)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = c.Repository.Transaction(func(tx *gorm.DB) error {
// 		if err := c.Repository.Delete(tx, dto, entity); err != nil {
// 			return err
// 		}

// 		_, err = c.VoidService.Delete(entity.GetFile())
// 		return err

// 	})

// 	return nil, err
// }

// func (c *ContextService) PutOrder(dto *request.AttachmentDto, body *req.OrderUpdateDto) (*response.AttachmentAdvancedResponseDto, error) {
// 	attachement, err := c.Repository.ValidateEntityExistence(dto)
// 	if err != nil || attachement.Order == body.Position {
// 		return common.NewResponse[response.AttachmentAdvancedResponseDto](attachement), err
// 	}

// 	err = c.Repository.Reorder(nil, attachement, body.Position)
// 	return common.NewResponse[response.AttachmentAdvancedResponseDto](attachement), err

// }

// func (c *ContextService) InitProjectFromTemplate(project *pre.ProjectEntity, readme bool) []entities.AttachmentEntity {
// 	var attachments []entities.AttachmentEntity

// 	if readme {
// 		entity := entities.NewAttachmentEntity()
// 		entity.Home = project.GetPath()
// 		entity.Create()

// 		entity.Name, entity.Path, entity.Size, entity.Type = "README.md", "/", 0, ".md"
// 		c.VoidService.Rename("/templates/readme.template.md", entity.GetFile(), true)
// 		attachments = append(attachments, *entity)
// 	}

// 	entity := entities.NewAttachmentEntity()
// 	entity.Home = project.GetPath()
// 	entity.Create()

// 	switch project.Type {
// 	case prt.P5js:
// 		entity.Name, entity.Path, entity.Size, entity.Type = "Main.js", "/", 0, ".js"
// 		c.VoidService.Rename("/templates/p5js.template.js", entity.GetFile(), true)

// 	case prt.Html:
// 		entity.Name, entity.Path, entity.Size, entity.Type = "index.html", "/", 0, ".html"
// 		c.VoidService.Rename("/templates/html.template.html", entity.GetFile(), true)

// 	case prt.Markdown:
// 		entity.Name, entity.Path, entity.Size, entity.Type = "index.md", "/", 0, ".md"
// 		c.VoidService.Rename("/templates/markdown.template.md", entity.GetFile(), true)

// 	// TODO:
// 	// case K3s:
// 	default:
// 		entity = nil
// 	}

// 	if entity != nil {
// 		attachments = append(attachments, *entity)
// 	}

// 	return attachments
// }
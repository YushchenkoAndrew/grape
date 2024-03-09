package project

import (
	"fmt"
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/project/dto/request"
	"grape/src/project/dto/response"
	e "grape/src/project/entities"
	repo "grape/src/project/repositories"
	"grape/src/project/types"

	"github.com/google/uuid"
)

type ProjectService struct {
	Repository *repo.ProjectRepositoryT

	// Link         i.Default[m.Link, m.LinkQueryDto]
	// File         i.Default[m.File, m.FileQueryDto]
	// Project      i.Default[m.Project, m.ProjectQueryDto]
	// Cron         *CronService
	// Subscription i.Default[m.Subscription, m.SubscribeQueryDto]
	// Metrics      i.Default[m.Metrics, m.MetricsQueryDto]
}

func NewProjectService(s *service.CommonService) *ProjectService {

	// (&repo.ProjectRepository{}).GetAll()
	return &ProjectService{
		Repository: repo.NewProjectRepository(s.DB),
		// Link:         NewLinkService(db, client),
		// File:         NewFileService(db, client),
		// Project:      NewProjectService(db, client),
		// Cron:         NewCronService(),
		// Subscription: NewSubscriptionService(db, client),
		// Metrics:      pods.NewMetricsService(db, client),
	}
}

func (c *ProjectService) FindOne(dto *request.ProjectDto) (*response.ProjectBasicResponseDto, error) {
	options := req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	})

	project, err := c.ValidateProjectExistence(options)
	if err != nil {
		return nil, err
	}

	return common.NewResponse[*response.ProjectBasicResponseDto](project), nil
}

func (c *ProjectService) FindAll(dto *request.ProjectDto) (common.PageResponseDto[[]response.ProjectBasicResponseDto], error) {
	options := req.NewRequest(dto, &request.ProjectDto{
		Statuses: []string{types.Active.String()},
	})

	total, projects := c.Repository.GetAllPage(*options)
	return common.PageResponseDto[[]response.ProjectBasicResponseDto]{
		Page:    options.Page,
		PerPage: options.Take,
		Total:   total,
		Result:  common.NewResponseMany[response.ProjectBasicResponseDto](projects),
	}, nil
}

func (c *ProjectService) Create(dto *request.ProjectDto, body *request.ProjectCreateDto) (*common.UuidResponseDto, error) {
	entity, err := c.Repository.Create(dto, body)
	if err != nil {
		return nil, err
	}

	return common.NewResponse[*common.UuidResponseDto](entity), nil
}

func (c *ProjectService) ValidateProjectExistence(dto *request.ProjectDto, relations ...repo.ProjectRelation) (*e.ProjectEntity, error) {
	if uuid.Validate(dto.ProjectIds[0]) != nil {
		return nil, fmt.Errorf("project id is invalid")
	}

	if project := c.Repository.GetOne(dto, relations...); project != nil {
		return project, nil
	}

	return nil, fmt.Errorf("project not found")
}

// 	// FIXME:
// 	// 	* Some strange error with page
// 	// 	* I guess the problim in limit value
// 	// 	* What I mean is what if we have several pages
// 	// 	* And First one is full, in this case I need to check another one
// 	// 	* To be sure that the new value was created
// 	// 	* I guess such functionality is missing !!!
// 	if res := c.db.Create(model); res.Error != nil {
// 		go logs.DefaultLog("/controllers/link.go", res.Error)
// 		return fmt.Errorf("Something unexpected happend: %v", res.Error)
// 	}

// 	c.precache(model, c.keys(model))
// 	return nil
// }

// func (c *ProjectService) Update(query *m.ProjectQueryDto, model *m.Project) ([]m.Project, error) {
// 	var res *gorm.DB
// 	client, suffix := c.query(query, c.db)

// 	var models = []m.Project{}
// 	if err, rows := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil || rows == 0 {
// 		return nil, fmt.Errorf("Requested model do not exist")
// 	}

// 	existed := model.Copy()
// 	for _, item := range models {
// 		if model.Name != "" && c.isExist(existed) && existed.ID != item.ID {
// 			return nil, fmt.Errorf("Requested project with name=%s has already existed", model.Name)
// 		}
// 	}

// 	client, _ = c.query(query, c.db.Model(&m.Project{}))
// 	if res = client.Updates(model); res.Error != nil {
// 		go logs.DefaultLog("/controllers/project.go", res.Error)
// 		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
// 	}

// 	for _, item := range models {
// 		// FIXME: !!!!
// 		// c.recache(&item, (existed.Name != "" && existed.ID == 0) || existed.Flag != model.Flag)
// 		c.recache(item.Copy().Fill(model), c.keys(&item), false)
// 		c.recache(item.Fill(model), c.keys(&item), false)
// 	}

// 	// Check if Name is not empty, if so that for some safety magers
// 	// lets replace this unique index with ID
// 	if query.IsOK(model) {
// 		return c.Read(query)
// 	}

// 	var result = []m.Project{}
// 	for _, item := range models {
// 		var model, err = c.Read(&m.ProjectQueryDto{ID: item.ID})
// 		if err != nil {
// 			return nil, err
// 		}

// 		result = append(result, model...)
// 	}
// 	return result, nil
// }

// func (c *ProjectService) Delete(query *m.ProjectQueryDto) (int, error) {
// 	var models []m.Project
// 	client, suffix := c.query(query, c.db)

// 	if err, _ := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil {
// 		return 0, err
// 	}

// 	for _, model := range models {
// 		c.recache(&model, c.keys(&model), true)
// 	}

// 	return len(models), client.Delete(&m.Project{}).Error
// }

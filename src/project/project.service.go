package project

import (
	"fmt"
	"grape/src/common/service"
	"grape/src/project/dto/request"
	e "grape/src/project/entities"
	repo "grape/src/project/repositories"

	"github.com/google/uuid"
)

type projectService struct {
	Repository repo.ProjectRepository
	// Link         i.Default[m.Link, m.LinkQueryDto]
	// File         i.Default[m.File, m.FileQueryDto]
	// Project      i.Default[m.Project, m.ProjectQueryDto]
	// Cron         *CronService
	// Subscription i.Default[m.Subscription, m.SubscribeQueryDto]
	// Metrics      i.Default[m.Metrics, m.MetricsQueryDto]
}

func NewProjectService(s *service.CommonService) *projectService {

	// (&repo.ProjectRepository{}).GetAll()
	return &projectService{
		Repository: *repo.NewProjectRepository(s.DB),
		// Link:         NewLinkService(db, client),
		// File:         NewFileService(db, client),
		// Project:      NewProjectService(db, client),
		// Cron:         NewCronService(),
		// Subscription: NewSubscriptionService(db, client),
		// Metrics:      pods.NewMetricsService(db, client),
	}
}

// type ProjectService struct {
// 	key string

// 	db     *gorm.DB
// 	client *redis.Client
// }

// func NewProjectService(db *gorm.DB, client *redis.Client) i.Default[m.Project, m.ProjectQueryDto] {
// 	return &ProjectService{key: "PROJECT", db: db, client: client}
// }

func (c *projectService) FindOne(project_id string) (interface{}, error) {
	// err, rows := helper.Getcache(c.db.Where("name = ?", model.Name), c.client, c.key, fmt.Sprintf("NAME=%s", model.Name), model)
	// return err != nil || rows != 0
	c.ValidateProjectExistence(project_id)

	return nil, nil
}

func (c *projectService) FindAll(dto *request.ProjectDto) (interface{}, error) {
	// err, rows := helper.Getcache(c.db.Where("name = ?", model.Name), c.client, c.key, fmt.Sprintf("NAME=%s", model.Name), model)
	// return err != nil || rows != 0

	return nil, nil
}

func (c *projectService) ValidateProjectExistence(project_id string, relations ...repo.ProjectRelation) (*e.ProjectEntity, error) {
	if uuid.Validate(project_id) != nil {
		return nil, fmt.Errorf("project id is invalid")
	}

	// TODO: Add dto creation
	if project := c.Repository.GetOne(nil, relations...); project != nil {
		return project, nil
	}

	return nil, fmt.Errorf("project not found")
}

// func (c *ProjectService) Create(model *m.Project) error {
// 	// Check if such project name exists
// 	if c.isExist(model) {
// 		return fmt.Errorf("Requested project with name=%s has already existed", model.Name)
// 	}

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

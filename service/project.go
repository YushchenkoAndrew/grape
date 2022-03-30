package service

import (
	"api/config"
	"api/helper"
	"api/logs"
	m "api/models"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type FullProjectService struct {
	Link    *LinkService
	File    *FileService
	Project *ProjectService
}

func NewFullProjectService(db *gorm.DB, client *redis.Client) *FullProjectService {
	return &FullProjectService{Link: NewLinkService(db, client), File: NewFileService(db, client), Project: NewProjectService(db, client)}
}

type ProjectService struct {
	key string

	db     *gorm.DB
	client *redis.Client
}

func NewProjectService(db *gorm.DB, client *redis.Client) *ProjectService {
	return &ProjectService{key: "PROJECT", db: db, client: client}
}

func (c *ProjectService) isExist(model *m.Project) bool {
	err, rows := helper.Getcache(c.db.Where("name = ?", model.Name), c.client, c.key, fmt.Sprintf("NAME=%s", model.Name), model)
	return err != nil || rows != 0
}

func (c *ProjectService) precache(model *m.Project) {
	helper.Precache(c.client, c.key, fmt.Sprintf("ID=%d", model.ID), model)
	helper.Precache(c.client, c.key, fmt.Sprintf("NAME=%s", model.Name), model)
}

func (c *ProjectService) recache(model *m.Project, delete bool) {
	helper.Delcache(c.client, c.key, fmt.Sprintf("ID=%d*", model.ID))
	helper.Delcache(c.client, c.key, fmt.Sprintf("NAME=%s*", model.Name))

	var keys = []string{fmt.Sprintf("FLAG=%s*", model.Flag), "CREATED_FROM=*", "CREATED_TO=*", "PAGE=*", "LIMIT=*"}
	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string) interface{} {
			if !strings.HasPrefix(str, "[") {
				var data m.Project
				json.Unmarshal([]byte(str), &data)
				if !delete {
					return *model
				}

				return nil
			}

			var data []m.Project
			var result []m.Project

			json.Unmarshal([]byte(str), &data)
			for _, item := range data {
				if item.Name != model.Name {
					result = append(result, item)
				} else if !delete {
					result = append(result, *model)
				}
			}

			if len(result) > 0 {
				return result
			}

			return nil

		})
	}
}

func (c *ProjectService) query(dto *m.ProjectQueryDto, client *gorm.DB) (*gorm.DB, string) {
	var suffix = ""

	if dto.ID > 0 {
		suffix += fmt.Sprintf("ID=%d", dto.ID)
		client = client.Where("id = ?", dto.ID)
	}

	if len(dto.Name) > 0 {
		suffix += fmt.Sprintf("NAME=%s", dto.Name)
		client = client.Where("name = ?", dto.Name)
	}

	if len(dto.Flag) > 0 {
		suffix += fmt.Sprintf("FLAG=%s", dto.Flag)
		client = client.Where("flag = ?", dto.Flag)
	}

	if !dto.CreatedFrom.IsZero() {
		suffix += fmt.Sprintf("CREATED_FROM=%s", dto.CreatedFrom.Format("2006-01-02"))
		client = client.Where("created_at >= ?", dto.CreatedFrom)
	}

	if !dto.CreatedTo.IsZero() {
		suffix += fmt.Sprintf("CREATED_TO=%s", dto.CreatedTo.Format("2006-01-02"))
		client = client.Where("created_at <= ?", dto.CreatedTo)
	}

	if dto.Page >= 0 {
		suffix += fmt.Sprintf("PAGE=%d", dto.Page)
		client = client.Offset(dto.Page * config.ENV.Items)
	}

	if dto.Limit > 0 {
		suffix += fmt.Sprintf("LIMIT=%d", dto.Page)
		client = client.Limit(dto.Limit)
	}

	return client, suffix
}

func (c *ProjectService) Create(model *m.Project) error {
	// Check if such project name exists
	if c.isExist(model) {
		return fmt.Errorf("Requested project with name=%s has already existed", model.Name)
	}

	if res := c.db.Create(model); res.Error != nil || res.RowsAffected == 0 {
		go logs.DefaultLog("/controllers/link.go", res.Error)
		return fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	c.precache(model)
	return nil
}

func (c *ProjectService) Read(query *m.ProjectQueryDto) ([]m.Project, error) {
	var model []m.Project
	client, suffix := c.query(query, c.db)

	err, _ := helper.Getcache(client.Order("created_at DESC"), c.client, c.key, suffix, &model)
	return model, err
}

func (c *ProjectService) Update(query *m.ProjectQueryDto, model *m.Project) ([]m.Project, error) {
	var res *gorm.DB
	client, suffix := c.query(query, c.db)

	var models = []m.Project{}
	if err, rows := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil || rows == 0 {
		return nil, fmt.Errorf("Requested model do not exist")
	}

	existed := model.Copy()
	for _, item := range models {
		if model.Name != "" && c.isExist(existed) && existed.ID != item.ID {
			return nil, fmt.Errorf("Requested project with name=%s has already existed", model.Name)
		}
	}

	client, _ = c.query(query, c.db.Model(&m.Project{}))
	if res = client.Updates(model); res.Error != nil || res.RowsAffected == 0 {
		go logs.DefaultLog("/controllers/project.go", res.Error)
		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	for _, item := range models {
		c.recache(&item, (existed.Name != "" && existed.ID == 0) || existed.Flag != model.Flag)
		c.recache(item.Fill(model), false)
	}

	// Check if Name is not empty, if so that for some safety magers
	// lets replace this unique index with ID
	if query.Name != "" {
		query.ID = models[0].ID
		query.Name = ""
	}

	return c.Read(query)
}

func (c *ProjectService) Delete(query *m.ProjectQueryDto) (int, error) {
	var models []m.Project
	client, suffix := c.query(query, c.db)

	if err, _ := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil {
		return 0, err
	}

	for _, model := range models {
		c.recache(&model, true)
	}

	return len(models), client.Delete(&m.Project{}).Error
}

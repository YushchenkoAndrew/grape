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

type LinkService struct {
	key string

	db     *gorm.DB
	client *redis.Client
}

func NewLinkService(db *gorm.DB, client *redis.Client) *LinkService {
	return &LinkService{key: "LINK", db: db, client: client}
}

func (c *LinkService) isExist(model *m.Link) bool {
	res := c.db.Where("project_id = ? AND name = ?", model.ProjectID, model.Name).Find(&model)
	return !(res.RowsAffected == 0)
}

func (c *LinkService) precache(model *m.Link) {
	helper.Precache(c.client, c.key, fmt.Sprintf("ID=%d", model.ID), model)
	helper.Precache(c.client, c.key, fmt.Sprintf("PROJECT_ID=%dNAME=%s", model.ProjectID, model.Name), model)
}

func (c *LinkService) recache(model *m.Link, delete bool) {
	helper.Delcache(c.client, c.key, fmt.Sprintf("ID=%d*", model.ID))

	var keys = []string{fmt.Sprintf("NAME=%s*", model.Name), fmt.Sprintf("PROJECT_ID=%d*", model.ProjectID), "PAGE=*", "LIMIT=*"}
	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string) interface{} {
			if !strings.HasPrefix(str, "[") {
				var data m.Link
				json.Unmarshal([]byte(str), &data)
				if !delete {
					return *model
				}

				return nil
			}

			var data []m.Link
			var result []m.Link

			json.Unmarshal([]byte(str), &data)
			for _, item := range data {
				if item.ID != model.ID {
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

func (c *LinkService) query(dto *m.LinkQueryDto, client *gorm.DB) (*gorm.DB, string) {
	var suffix = ""

	if dto.ID > 0 {
		suffix += fmt.Sprintf("ID=%d", dto.ID)
		client = client.Where("id = ?", dto.ID)
	}

	if dto.ProjectID > 0 {
		suffix += fmt.Sprintf("PROJECT_ID=%d", dto.ProjectID)
		client = client.Where("project_id = ?", dto.ProjectID)
	}

	if len(dto.Name) > 0 {
		suffix += fmt.Sprintf("NAME=%s", dto.Name)
		client = client.Where("name = ?", dto.Name)
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

func (c *LinkService) Create(model *m.Link) error {
	// Check if such project_id exists
	if err, rows := helper.Getcache(c.db.Where("id = ?", model.ProjectID), c.client, "PROJECT", fmt.Sprintf("ID=%d", model.ProjectID), &[]m.Project{}); err != nil || rows == 0 {
		return fmt.Errorf("Requested project_id=%d do not exist", model.ProjectID)
	}

	var res *gorm.DB
	var existed = model.Copy()

	if c.isExist(existed) {
		if res = c.db.Model(&m.Link{}).Where("id = ?", existed.ID).Updates(model); res.Error == nil {
			c.recache(existed.Fill(model), false)
		}
	} else if res = c.db.Create(model); res.Error == nil {
		c.precache(model)
	}

	if res.Error != nil || res.RowsAffected == 0 {
		go logs.DefaultLog("/controllers/link.go", res.Error)
		return fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	return nil
}

func (c *LinkService) Read(query *m.LinkQueryDto) ([]m.Link, error) {
	var model []m.Link
	client, suffix := c.query(query, c.db)

	err, _ := helper.Getcache(client.Order("updated_at DESC"), c.client, c.key, suffix, &model)
	return model, err
}

func (c *LinkService) Update(query *m.LinkQueryDto, model *m.Link) ([]m.Link, error) {
	var res *gorm.DB
	client, suffix := c.query(query, c.db)

	var models = []m.Link{}
	if err, rows := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil || rows == 0 {
		return nil, fmt.Errorf("Requested id=%d do not exist", model.ID)
	}

	client, _ = c.query(query, c.db.Model(&m.Link{}))
	if res = client.Updates(model); res.Error != nil || res.RowsAffected == 0 {
		go logs.DefaultLog("/controllers/link.go", res.Error)
		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	for _, existed := range models {
		c.recache(existed.Fill(model), false)
	}

	return c.Read(query)
}

func (c *LinkService) Delete(query *m.LinkQueryDto) (int, error) {
	var models []m.Link
	client, suffix := c.query(query, c.db)

	if err, _ := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil {
		return 0, err
	}

	for _, model := range models {
		c.recache(&model, true)
	}

	return len(models), client.Delete(&m.Link{}).Error
}

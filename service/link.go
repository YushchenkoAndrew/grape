package service

import (
	"api/config"
	"api/helper"
	m "api/models"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type LinkService struct {
	db     *gorm.DB
	client *redis.Client
}

const KEY = "LINK"

func NewLinkService(db *gorm.DB, client *redis.Client) *LinkService {
	return &LinkService{db, client}
}

func (c *LinkService) isExist(model *m.Link) bool {
	res := c.db.Where("project_id = ? AND name = ?", model.ProjectID, model.Name).Find(&model)
	return !(res.RowsAffected == 0)
}

func (c *LinkService) precache(model *m.Link) {
	helper.Precache(c.client, KEY, fmt.Sprintf("ID=%d", model.ID), model)
	helper.Precache(c.client, KEY, fmt.Sprintf("PROJECT_ID=%dNAME=%s", model.ProjectID, model.Name), model)
}

func (c *LinkService) delcache(model *m.Link) {
	helper.Delcache(c.client, KEY, fmt.Sprintf("ID=%d", model.ID), model)
	helper.Delcache(c.client, KEY, fmt.Sprintf("PROJECT_ID=%dNAME=%s", model.ProjectID, model.Name), model)
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

	if dto.Name != "" {
		suffix += fmt.Sprintf("NAME=%s", dto.Name)
		client = client.Where("name = ?", dto.Name)
	}

	if dto.Page >= 0 {
		suffix += fmt.Sprintf("PAGE=%d", dto.Page)
		client = client.Offset(dto.Page * config.ENV.Items)
	}

	if dto.Limit >= 0 {
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
		c.delcache(model)
		res = c.db.Model(&m.Link{}).Where("id = ?", existed.ID).Updates(model)
	} else {
		res = c.db.Create(model)
	}

	if res.Error != nil || res.RowsAffected == 0 {
		// go logs.DefaultLog("/controllers/link.go", res.Error)
		return fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	c.precache(model)
	return nil
}

func (c *LinkService) Read(query *m.LinkQueryDto) ([]m.Link, error) {
	var model []m.Link
	client, suffix := c.query(query, c.db)

	err, _ := helper.Getcache(client.Order("updated_at DESC"), c.client, KEY, suffix, &model)
	return model, err
}

func (c *LinkService) Update(query *m.LinkQueryDto, model *m.Link) error {
	var res *gorm.DB
	client, suffix := c.query(query, c.db)

	if err, rows := helper.Getcache(client, c.client, KEY, suffix, &[]m.Link{}); err != nil || rows == 0 {
		return fmt.Errorf("Requested id=%d do not exist", model.ID)
	} else {
		client, _ := c.query(query, c.db.Model(&m.Link{}))

		c.delcache(model)
		res = client.Updates(model)
	}

	if res.Error != nil || res.RowsAffected == 0 {
		// go logs.DefaultLog("/controllers/link.go", res.Error)
		return fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	c.precache(model)
	return nil
}

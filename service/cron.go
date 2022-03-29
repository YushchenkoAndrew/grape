package service

import (
	"api/config"
	"api/helper"
	"api/logs"
	m "api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CronService struct {
}

func NewCronService(db *gorm.DB, client *redis.Client) *CronService {
	return &CronService{}
}

func (c *CronService) Create(dto *m.CronCreateDto) (*m.CronEntity, error) {
	var body, err = json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	var salt, token = helper.BotToken()
	if req, err = http.NewRequest("POST", fmt.Sprintf("%s/cron/subscribe?key=%s", config.ENV.BotUrl, token), bytes.NewBuffer(body)); err != nil {
		return nil, err
	}

	req.Header.Set("X-Custom-Header", salt)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Bot request error")
	}

	var model m.CronEntity
	if err = json.NewDecoder(res.Body).Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *CronService) Read(query *m.FileQueryDto) ([]m.File, error) {
	var model []m.File
	client, suffix := c.query(query, c.db)

	err, _ := helper.Getcache(client.Order("updated_at DESC"), c.client, c.key, suffix, &model)
	return model, err
}

func (c *CronService) Update(query *m.FileQueryDto, model *m.File) ([]m.File, error) {
	var res *gorm.DB
	client, suffix := c.query(query, c.db)

	var models = []m.File{}
	if err, rows := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil || rows == 0 {
		return nil, fmt.Errorf("Requested model do not exist")
	}

	client, _ = c.query(query, c.db.Model(&m.File{}))
	if res = client.Updates(model); res.Error != nil || res.RowsAffected == 0 {
		go logs.DefaultLog("/controllers/file.go", res.Error)
		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	for _, existed := range models {
		c.recache(existed.Fill(model), false)
	}

	return c.Read(query)
}

func (c *CronService) Delete(query *m.FileQueryDto) (int, error) {
	var models []m.File
	client, suffix := c.query(query, c.db)

	if err, _ := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil {
		return 0, err
	}

	for _, model := range models {
		c.recache(&model, true)
	}

	return len(models), client.Delete(&m.File{}).Error
}

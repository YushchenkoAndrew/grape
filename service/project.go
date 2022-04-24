package service

import (
	"api/config"
	"api/helper"
	i "api/interfaces/service"
	"api/logs"
	m "api/models"
	"api/service/k3s/pods"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type FullProjectService struct {
	Link         i.Default[m.Link, m.LinkQueryDto]
	File         i.Default[m.File, m.FileQueryDto]
	Project      i.Default[m.Project, m.ProjectQueryDto]
	Cron         *CronService
	Subscription i.Default[m.Subscription, m.SubscribeQueryDto]
	Metrics      i.Default[m.Metrics, m.MetricsQueryDto]
}

func NewFullProjectService(db *gorm.DB, client *redis.Client) *FullProjectService {
	return &FullProjectService{
		Link:         NewLinkService(db, client),
		File:         NewFileService(db, client),
		Project:      NewProjectService(db, client),
		Cron:         NewCronService(),
		Subscription: NewSubscriptionService(db, client),
		Metrics:      pods.NewMetricsService(db, client),
	}
}

type ProjectService struct {
	key string

	db     *gorm.DB
	client *redis.Client
}

func NewProjectService(db *gorm.DB, client *redis.Client) i.Default[m.Project, m.ProjectQueryDto] {
	return &ProjectService{key: "PROJECT", db: db, client: client}
}

func (c *ProjectService) keys(model *m.Project) []string {
	return []string{fmt.Sprintf("FLAG=%s*", model.Flag), "", "CREATED_FROM=*", "CREATED_TO=*", "PAGE=*", "LIMIT=*"}
}

func (c *ProjectService) isExist(model *m.Project) bool {
	err, rows := helper.Getcache(c.db.Where("name = ?", model.Name), c.client, c.key, fmt.Sprintf("NAME=%s", model.Name), model)
	return err != nil || rows != 0
}

func (c *ProjectService) precache(model *m.Project, keys []string) {
	helper.Precache(c.client, c.key, fmt.Sprintf("ID=%d", model.ID), model)
	helper.Precache(c.client, c.key, fmt.Sprintf("NAME=%s", model.Name), model)

	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string, k string) interface{} {
			if !strings.HasPrefix(str, "[") {
				return *model
			}

			var data []m.Project
			json.Unmarshal([]byte(str), &data)
			for _, item := range strings.Split(k, "#") {
				var res = strings.Split(item, "=")

				switch res[0] {
				case "ID":
					if id, _ := strconv.Atoi(res[1]); model.ID != uint32(id) {
						return data
					}

				case "NAME":
					if model.Name != res[1] {
						return data
					}

				case "FLAG":
					if model.Flag != res[1] {
						return data
					}

				case "CREATED_FROM":
					if created_from, _ := time.Parse("2006-01-02", res[1]); created_from.Year() < model.CreatedAt.Year() || created_from.YearDay() < model.CreatedAt.YearDay() {
						return data
					}

				case "CREATED_TO":
					if created_to, _ := time.Parse("2006-01-02", res[1]); created_to.Year() > model.CreatedAt.Year() || created_to.YearDay() > model.CreatedAt.YearDay() {
						return data
					}

				case "LIMIT":
					if limit, _ := strconv.Atoi(res[1]); limit <= len(data) {
						return data
					}
				}
			}

			return append(data, *model)
		})
	}
}

func (c *ProjectService) deepcache(models []m.Project, key string) interface{} {
	var suffix []string
	for _, item := range strings.Split(key, "#") {
		var res = strings.Split(item, "=")

		switch res[0] {
		case "LIMIT":
			if limit, _ := strconv.Atoi(res[1]); limit != len(models)+1 {
				if len(models) != 0 {
					return models
				}
				return nil
			}

		case "PAGE":
			page, _ := strconv.Atoi(res[1])
			suffix = append(suffix, fmt.Sprintf("PAGE=%d", page+1))
			continue
		}

		suffix = append(suffix, item)
	}

	var items []m.Project
	if err := helper.Popcache(c.client, c.key, strings.Join(suffix, "#"), &items); err == nil {
		if len(items) < 1 {
			return nil
		}

		go c.deepcache(items[1:], strings.Join(suffix, "#"))
		return append(models, items[0])
	}

	if len(models) != 0 {
		return models
	}
	return nil
}

func (c *ProjectService) postfilter(data []m.Project, suffix string) []m.Project {
	var result = []m.Project{}

ITEM:
	for _, item := range data {
		for _, key := range strings.Split(suffix, "#") {
			var res = strings.Split(key, "=")

			switch res[0] {
			case "ID":
				if id, _ := strconv.Atoi(res[1]); item.ID != uint32(id) {
					continue ITEM
				}

			case "NAME":
				if item.Name != res[1] {
					continue ITEM
				}

			case "FLAG":
				if item.Flag != res[1] {
					continue ITEM
				}

			case "CREATED_FROM":
				if created_from, _ := time.Parse("2006-01-02", res[1]); created_from.Year() < item.CreatedAt.Year() || created_from.YearDay() < item.CreatedAt.YearDay() {
					continue ITEM
				}

			case "CREATED_TO":
				if created_to, _ := time.Parse("2006-01-02", res[1]); created_to.Year() > item.CreatedAt.Year() || created_to.YearDay() > item.CreatedAt.YearDay() {
					continue ITEM
				}
			}
		}

		result = append(result, item)
	}

	return result
}

func (c *ProjectService) recache(model *m.Project, keys []string, delete bool) {
	helper.Delcache(c.client, c.key, fmt.Sprintf("ID=%d*", model.ID))
	helper.Delcache(c.client, c.key, fmt.Sprintf("NAME=%s*", model.Name))

	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string, suffix string) interface{} {
			if !strings.HasPrefix(str, "[") {
				str = fmt.Sprintf("[%s]", str)
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

			// Postfilter elements with cache query
			result = c.postfilter(result, suffix)

			// Check if size of an array was changed
			if delete {
				return c.deepcache(result, suffix)
			}

			return result
		})
	}
}

func (c *ProjectService) query(dto *m.ProjectQueryDto, client *gorm.DB) (*gorm.DB, string) {
	var suffix []string

	if dto.ID > 0 {
		suffix = append(suffix, fmt.Sprintf("ID=%d", dto.ID))
		client = client.Where("id = ?", dto.ID)
	}

	if len(dto.Name) > 0 {
		suffix = append(suffix, fmt.Sprintf("NAME=%s", dto.Name))
		client = client.Where("name = ?", dto.Name)
	}

	if len(dto.Flag) > 0 {
		suffix = append(suffix, fmt.Sprintf("FLAG=%s", dto.Flag))
		client = client.Where("flag = ?", dto.Flag)
	}

	if !dto.CreatedFrom.IsZero() {
		suffix = append(suffix, fmt.Sprintf("CREATED_FROM=%s", dto.CreatedFrom.Format("2006-01-02")))
		client = client.Where("created_at >= ?", dto.CreatedFrom)
	}

	if !dto.CreatedTo.IsZero() {
		suffix = append(suffix, fmt.Sprintf("CREATED_TO=%s", dto.CreatedTo.Format("2006-01-02")))
		client = client.Where("created_at <= ?", dto.CreatedTo)
	}

	if dto.Page >= 0 {
		var limit = config.ENV.Items
		if dto.Limit > 0 {
			limit = dto.Limit
		}

		suffix = append(suffix, fmt.Sprintf("PAGE=%d", dto.Page))
		client = client.Offset(dto.Page * limit)

		suffix = append(suffix, fmt.Sprintf("LIMIT=%d", limit))
		client = client.Limit(limit)

	}

	return client, strings.Join(suffix, "#")
}

func (c *ProjectService) Create(model *m.Project) error {
	// Check if such project name exists
	if c.isExist(model) {
		return fmt.Errorf("Requested project with name=%s has already existed", model.Name)
	}

	// FIXME:
	// 	* Some strange error with page
	// 	* I guess the problim in limit value
	// 	* What I mean is what if we have several pages
	// 	* And First one is full, in this case I need to check another one
	// 	* To be sure that the new value was created
	// 	* I guess such functionality is missing !!!
	if res := c.db.Create(model); res.Error != nil {
		go logs.DefaultLog("/controllers/link.go", res.Error)
		return fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	c.precache(model, c.keys(model))
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
	if res = client.Updates(model); res.Error != nil {
		go logs.DefaultLog("/controllers/project.go", res.Error)
		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	for _, item := range models {
		// FIXME: !!!!
		// c.recache(&item, (existed.Name != "" && existed.ID == 0) || existed.Flag != model.Flag)
		c.recache(item.Copy().Fill(model), c.keys(&item), false)
		c.recache(item.Fill(model), c.keys(&item), false)
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
		c.recache(&model, c.keys(&model), true)
	}

	return len(models), client.Delete(&m.Project{}).Error
}

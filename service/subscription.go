package service

import (
	"api/config"
	"api/helper"
	"api/logs"
	m "api/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type FullSubscriptionService struct {
	Cron         *CronService
	Subscription *SubscriptionService
}

func NewFullSubscriptionService(db *gorm.DB, client *redis.Client) *FullSubscriptionService {
	return &FullSubscriptionService{Cron: NewCronService(), Subscription: NewSubscriptionService(db, client)}
}

type SubscriptionService struct {
	key string

	db     *gorm.DB
	client *redis.Client
}

func NewSubscriptionService(db *gorm.DB, client *redis.Client) *SubscriptionService {
	return &SubscriptionService{key: "SUBSCRIPTION", db: db, client: client}
}

func (c *SubscriptionService) isExist(model *m.Subscription) bool {
	err, rows := helper.Getcache(c.db.Where("cron_id = ?", model.CronID), c.client, c.key, fmt.Sprintf("CRON_ID=%s", model.CronID), model)
	return err != nil || rows != 0
}

func (c *SubscriptionService) precache(model *m.Subscription) {
	helper.Precache(c.client, c.key, fmt.Sprintf("ID=%d", model.ID), model)
	helper.Precache(c.client, c.key, fmt.Sprintf("TOKEN=%s", model.Token), model)
	helper.Precache(c.client, c.key, fmt.Sprintf("CRON_ID=%s", model.CronID), model)

	var keys = []string{fmt.Sprintf("NAME=%s*", model.Name), fmt.Sprintf("PROJECT_ID=%d*", model.ProjectID), "PAGE=*", "LIMIT=*"}
	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string, k string) interface{} {
			var data []m.Subscription
			if !strings.HasPrefix(str, "[") {
				data = make([]m.Subscription, 1)
				json.Unmarshal([]byte(str), &data[0])
			} else {
				json.Unmarshal([]byte(str), &data)
			}

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

				case "CRON_ID":
					if model.CronID != res[1] {
						return data
					}

				case "PROJECT_ID":
					if id, _ := strconv.Atoi(res[1]); model.ProjectID != uint32(id) {
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

func (c *SubscriptionService) deepcache(models []m.Subscription, key string) interface{} {
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

	var items []m.Subscription
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

func (c *SubscriptionService) recache(model *m.Subscription, delete bool) {
	helper.Delcache(c.client, c.key, fmt.Sprintf("ID=%d*", model.ID))
	helper.Delcache(c.client, c.key, fmt.Sprintf("TOKEN=%s*", model.Token))
	helper.Delcache(c.client, c.key, fmt.Sprintf("CRON_ID=%s*", model.CronID))

	var keys = []string{fmt.Sprintf("NAME=%s*", model.Name), fmt.Sprintf("PROJECT_ID=%d*", model.ProjectID), "PAGE=*", "LIMIT=*"}
	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string, suffix string) interface{} {
			if !strings.HasPrefix(str, "[") {
				var data m.Subscription
				json.Unmarshal([]byte(str), &data)
				if !delete {
					return *model
				}

				return nil
			}

			var data []m.Subscription
			var result []m.Subscription

			json.Unmarshal([]byte(str), &data)
			for _, item := range data {
				if item.CronID != model.CronID {
					result = append(result, item)
				} else if !delete {
					result = append(result, *model)
				}
			}

			// Check if size of an array was changed
			if delete {
				return c.deepcache(result, suffix)
			}

			return result
		})
	}
}

func (c *SubscriptionService) query(dto *m.SubscribeQueryDto, client *gorm.DB) (*gorm.DB, string) {
	var suffix []string

	if dto.ID > 0 {
		suffix = append(suffix, fmt.Sprintf("ID=%d", dto.ID))
		client = client.Where("id = ?", dto.ID)
	}

	if dto.ProjectID > 0 {
		suffix = append(suffix, fmt.Sprintf("PROJECT_ID=%d", dto.ProjectID))
		client = client.Where("project_id = ?", dto.ProjectID)
	}

	if len(dto.Name) > 0 {
		suffix = append(suffix, fmt.Sprintf("NAME=%s", dto.Name))
		client = client.Where("name = ?", dto.Name)
	}

	if len(dto.CronID) > 0 {
		suffix = append(suffix, fmt.Sprintf("CRON_ID=%s", dto.CronID))
		client = client.Where("cron_id = ?", dto.CronID)
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

func (c *SubscriptionService) Create(model *m.Subscription) error {
	// Check if such project_id exists
	if err, rows := helper.Getcache(c.db.Where("id = ?", model.ProjectID), c.client, "PROJECT", fmt.Sprintf("ID=%d", model.ProjectID), &[]m.Project{}); err != nil || rows == 0 {
		return fmt.Errorf("Requested project_id=%d do not exist", model.ProjectID)
	}

	var res *gorm.DB
	var existed = model.Copy()

	if c.isExist(existed) {
		if res = c.db.Model(&m.Subscription{}).Where("id = ?", existed.ID).Updates(model); res.Error == nil {
			c.recache(existed.Fill(model), false)
		}
	} else if res = c.db.Create(model); res.Error == nil {
		c.precache(model)
	}

	if res.Error != nil {
		go logs.DefaultLog("/controllers/subscription.go", res.Error)
		return fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	return nil
}

func (c *SubscriptionService) Read(query *m.SubscribeQueryDto) ([]m.Subscription, error) {
	var model []m.Subscription
	client, suffix := c.query(query, c.db)

	err, _ := helper.Getcache(client.Order("created_at DESC"), c.client, c.key, suffix, &model)
	return model, err
}

func (c *SubscriptionService) Update(query *m.SubscribeQueryDto, model *m.Subscription) ([]m.Subscription, error) {
	var res *gorm.DB
	client, suffix := c.query(query, c.db)

	var models = []m.Subscription{}
	if err, rows := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil || rows == 0 {
		return nil, fmt.Errorf("Requested model do not exist")
	}

	existed := model.Copy()
	for _, item := range models {
		if model.CronID != "" && c.isExist(existed) && existed.ID != item.ID {
			return nil, fmt.Errorf("Requested project with name=%s has already existed", model.Name)
		}
	}

	client, _ = c.query(query, c.db.Model(&m.Subscription{}))
	if res = client.Updates(model); res.Error != nil {
		go logs.DefaultLog("/controllers/subscription.go", res.Error)
		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	for _, item := range models {
		c.recache(&item, existed.CronID != "" && existed.ID == 0)
		c.recache(item.Fill(model), false)
	}

	// Check if Name is not empty, if so that for some safety magers
	// lets replace this unique index with ID
	if query.CronID != "" {
		query.ID = models[0].ID
		query.CronID = ""
	}

	return c.Read(query)
}

func (c *SubscriptionService) Delete(query *m.SubscribeQueryDto) (int, error) {
	var models []m.Subscription
	client, suffix := c.query(query, c.db)

	if err, _ := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil {
		return 0, err
	}

	for _, model := range models {
		c.recache(&model, true)
	}

	return len(models), client.Delete(&m.Subscription{}).Error
}

package service

import (
	"api/config"
	"api/helper"
	i "api/interfaces/service"
	"api/logs"
	m "api/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type LinkService struct {
	key string

	db     *gorm.DB
	client *redis.Client
}

func NewLinkService(db *gorm.DB, client *redis.Client) i.Default[m.Link, m.LinkQueryDto] {
	return &LinkService{key: "LINK", db: db, client: client}
}

func (c *LinkService) keys(model *m.Link) []string {
	return []string{fmt.Sprintf("NAME=%s*", model.Name), fmt.Sprintf("PROJECT_ID=%d*", model.ProjectID), "", "PAGE=*", "LIMIT=*"}
}

func (c *LinkService) isExist(model *m.Link) bool {
	var res = []*m.Link{model}
	err, rows := helper.Getcache(c.db.Where("project_id = ? AND name = ?", model.ProjectID, model.Name), c.client, c.key, fmt.Sprintf("PROJECT_ID=%d#NAME=%s", model.ProjectID, model.Name), &res)
	return err != nil || rows != 0
}

func (c *LinkService) precache(model *m.Link, keys []string) {
	helper.Precache(c.client, c.key, fmt.Sprintf("ID=%d", model.ID), model)
	helper.Precache(c.client, c.key, fmt.Sprintf("PROJECT_ID=%d#NAME=%s", model.ProjectID, model.Name), model)

	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string, k string) interface{} {
			var data []m.Link
			if !strings.HasPrefix(str, "[") {
				data = make([]m.Link, 1)
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

func (c *LinkService) deepcache(models []m.Link, key string) interface{} {
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

	var items []m.Link
	if err := helper.Popcache(c.client, c.key, strings.Join(suffix, "#"), &items); err == nil {
		if len(items) == 0 {
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

func (c *LinkService) postfilter(data []m.Link, suffix string) []m.Link {
	var result = []m.Link{}

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

			case "PROJECT_ID":
				if id, _ := strconv.Atoi(res[1]); item.ProjectID != uint32(id) {
					continue ITEM
				}
			}
		}

		result = append(result, item)
	}

	return result
}

func (c *LinkService) recache(model *m.Link, keys []string, delete bool) {
	helper.Delcache(c.client, c.key, fmt.Sprintf("ID=%d*", model.ID))

	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string, suffix string) interface{} {
			if !strings.HasPrefix(str, "[") {
				str = fmt.Sprintf("[%s]", str)
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

func (c *LinkService) query(dto *m.LinkQueryDto, client *gorm.DB) (*gorm.DB, string) {
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

func (c *LinkService) Create(model *m.Link) error {
	// Check if such project_id exists
	if err, rows := helper.Getcache(c.db.Where("id = ?", model.ProjectID), c.client, "PROJECT", fmt.Sprintf("ID=%d", model.ProjectID), &[]m.Project{}); err != nil || rows == 0 {
		return fmt.Errorf("Requested project_id=%d do not exist", model.ProjectID)
	}

	var res *gorm.DB
	var existed = model.Copy()

	if c.isExist(existed) {
		if res = c.db.Model(&m.Link{}).Where("id = ?", existed.ID).Updates(model); res.Error == nil {
			c.recache(existed.Copy().Fill(model), c.keys(existed), false)
			c.recache(existed.Fill(model), c.keys(existed), false)
		}
	} else if res = c.db.Create(model); res.Error == nil {
		c.precache(model, c.keys(model))
	}

	if res.Error != nil {
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
		return nil, fmt.Errorf("Requested model do not exist")
	}

	client, _ = c.query(query, c.db.Model(&m.Link{}))
	if res = client.Updates(model); res.Error != nil {
		go logs.DefaultLog("/controllers/link.go", res.Error)
		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	for _, existed := range models {
		c.recache(existed.Copy().Fill(model), c.keys(&existed), false)
		c.recache(existed.Fill(model), c.keys(&existed), false)
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
		c.recache(&model, c.keys(&model), true)
	}

	return len(models), client.Delete(&m.Link{}).Error
}

package pattern

import "grape/src/common/service"

// import (
// 	"grape/config"
// 	"grape/helper"
// 	i "grape/interfaces/service"
// 	"grape/logs"
// 	m "grape/models"
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/go-redis/redis/v8"
// 	"gorm.io/gorm"
// )

type PatternService struct{}

func NewPatternService(s *service.CommonService) *PatternService {
	return &PatternService{}
}

// func (c *PatternService) keys(model *m.Pattern) []string {
// 	return []string{fmt.Sprintf("MODE=%s*", model.Mode), fmt.Sprintf("COLORS=%d*", model.Colors), "", "PAGE=*", "LIMIT=*"}
// }

// func (c *PatternService) isExist(model *m.Pattern) bool {
// 	var res = []*m.Pattern{model}
// 	result := c.db.Where("path = ?", model.Path).Find(&res)
// 	return result.Error != nil || result.RowsAffected != 0
// }

// func (c *PatternService) precache(model *m.Pattern, keys []string) {
// 	helper.Precache(c.client, c.key, fmt.Sprintf("ID=%d", model.ID), model)

// 	for _, key := range keys {
// 		helper.Recache(c.client, c.key, key, func(str string, k string) interface{} {
// 			var data []m.Pattern
// 			if !strings.HasPrefix(str, "[") {
// 				data = make([]m.Pattern, 1)
// 				json.Unmarshal([]byte(str), &data[0])
// 			} else {
// 				json.Unmarshal([]byte(str), &data)
// 			}

// 			for _, item := range strings.Split(k, "#") {
// 				var res = strings.Split(item, "=")

// 				switch res[0] {
// 				case "ID":
// 					if id, _ := strconv.Atoi(res[1]); model.ID != uint32(id) {
// 						return data
// 					}

// 				case "MODE":
// 					if model.Mode != res[1] {
// 						return data
// 					}

// 				case "COLORS":
// 					if colors, _ := strconv.Atoi(res[1]); uint32(model.Colors) != uint32(colors) {
// 						return data
// 					}

// 				case "LIMIT":
// 					if limit, _ := strconv.Atoi(res[1]); limit <= len(data) {
// 						return data
// 					}
// 				}
// 			}

// 			return append(data, *model)
// 		})
// 	}
// }

// func (c *PatternService) deepcache(models []m.Pattern, key string) interface{} {
// 	var suffix []string
// 	for _, item := range strings.Split(key, "#") {
// 		var res = strings.Split(item, "=")

// 		switch res[0] {
// 		case "LIMIT":
// 			if limit, _ := strconv.Atoi(res[1]); limit != len(models)+1 {
// 				if len(models) != 0 {
// 					return models
// 				}
// 				return nil
// 			}

// 		case "PAGE":
// 			page, _ := strconv.Atoi(res[1])
// 			suffix = append(suffix, fmt.Sprintf("PAGE=%d", page+1))
// 			continue
// 		}

// 		suffix = append(suffix, item)
// 	}

// 	var items []m.Pattern
// 	if err := helper.Popcache(c.client, c.key, strings.Join(suffix, "#"), &items); err == nil {
// 		if len(items) == 0 {
// 			return nil
// 		}

// 		go c.deepcache(items[1:], strings.Join(suffix, "#"))
// 		return append(models, items[0])
// 	}

// 	if len(models) != 0 {
// 		return models
// 	}
// 	return nil
// }

// func (c *PatternService) postfilter(data []m.Pattern, suffix string) []m.Pattern {
// 	var result = []m.Pattern{}

// ITEM:
// 	for _, item := range data {
// 		for _, key := range strings.Split(suffix, "#") {
// 			var res = strings.Split(key, "=")

// 			switch res[0] {
// 			case "ID":
// 				if id, _ := strconv.Atoi(res[1]); item.ID != uint32(id) {
// 					continue ITEM
// 				}

// 			case "MODE":
// 				if item.Mode != res[1] {
// 					continue ITEM
// 				}

// 			case "COLORS":
// 				if colors, _ := strconv.Atoi(res[1]); uint32(item.Colors) != uint32(colors) {
// 					continue ITEM
// 				}
// 			}
// 		}

// 		result = append(result, item)
// 	}

// 	return result
// }

// func (c *PatternService) recache(model *m.Pattern, keys []string, delete bool) {
// 	helper.Delcache(c.client, c.key, fmt.Sprintf("ID=%d*", model.ID))

// 	for _, key := range keys {
// 		helper.Recache(c.client, c.key, key, func(str string, suffix string) interface{} {
// 			if !strings.HasPrefix(str, "[") {
// 				str = fmt.Sprintf("[%s]", str)
// 			}

// 			var data []m.Pattern
// 			var result []m.Pattern

// 			json.Unmarshal([]byte(str), &data)
// 			for _, item := range data {
// 				if item.ID != model.ID {
// 					result = append(result, item)
// 				} else if !delete {
// 					result = append(result, *model)
// 				}
// 			}

// 			// Postfilter elements with cache query
// 			result = c.postfilter(result, suffix)

// 			// Check if size of an array was changed
// 			if delete {
// 				return c.deepcache(result, suffix)
// 			}

// 			if len(result) != 0 {
// 				return result
// 			}

// 			return nil
// 		})
// 	}
// }

// func (c *PatternService) query(dto *m.PatternQueryDto, client *gorm.DB) (*gorm.DB, string) {
// 	var suffix []string

// 	if dto.ID > 0 {
// 		suffix = append(suffix, fmt.Sprintf("ID=%d", dto.ID))
// 		client = client.Where("id = ?", dto.ID)
// 	}

// 	if len(dto.Mode) > 0 {
// 		suffix = append(suffix, fmt.Sprintf("MODE=%s", dto.Mode))
// 		client = client.Where("mode = ?", dto.Mode)
// 	}

// 	if dto.Colors > 0 {
// 		suffix = append(suffix, fmt.Sprintf("COLORS=%d", dto.Colors))
// 		client = client.Where("colors = ?", dto.Colors)
// 	}

// 	if dto.Page >= 0 {
// 		var limit = config.ENV.Items
// 		if dto.Limit > 0 {
// 			limit = dto.Limit
// 		}

// 		suffix = append(suffix, fmt.Sprintf("PAGE=%d", dto.Page))
// 		client = client.Offset(dto.Page * limit)

// 		suffix = append(suffix, fmt.Sprintf("LIMIT=%d", limit))
// 		client = client.Limit(limit)
// 	}

// 	return client, strings.Join(suffix, "#")
// }

// func (c *PatternService) Create(model *m.Pattern) error {
// 	// Check if such pattern path already exists
// 	if c.isExist(model) {
// 		return fmt.Errorf("Requested pattern with path=%s already exists", model.Path)
// 	}

// 	if res := c.db.Create(model); res.Error != nil {
// 		go logs.DefaultLog("/service/pattern.go", res.Error)
// 		return fmt.Errorf("Something unexpected happend: %v", res.Error)
// 	}

// 	c.precache(model, c.keys(model))
// 	return nil
// }

// func (c *PatternService) Read(query *m.PatternQueryDto) ([]m.Pattern, error) {
// 	var model []m.Pattern
// 	client, suffix := c.query(query, c.db)

// 	err, _ := helper.Getcache(client.Order("created_at DESC"), c.client, c.key, suffix, &model)
// 	return model, err
// }

// func (c *PatternService) Update(query *m.PatternQueryDto, model *m.Pattern) ([]m.Pattern, error) {
// 	var res *gorm.DB
// 	client, suffix := c.query(query, c.db)

// 	var models = []m.Pattern{}
// 	if err, rows := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil || rows == 0 {
// 		return nil, fmt.Errorf("Requested model do not exist")
// 	}

// 	existed := model.Copy()
// 	for _, item := range models {
// 		if model.Path != "" && c.isExist(existed) && existed.ID != item.ID {
// 			return nil, fmt.Errorf("Requested pattern with path='%s' has already existed", model.Path)
// 		}
// 	}

// 	client, _ = c.query(query, c.db.Model(&m.Pattern{}))
// 	if res = client.Updates(model); res.Error != nil {
// 		go logs.DefaultLog("/service/pattern.go", res.Error)
// 		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
// 	}

// 	for _, item := range models {
// 		// FIXME: CHECK WITH PROJECT !!!
// 		// c.recache(&item, (existed.Path != "" && existed.ID == 0))
// 		c.recache(item.Copy().Fill(model), c.keys(&item), false)
// 		c.recache(item.Fill(model), c.keys(&item), false)
// 	}

// 	if query.IsOK(model) {
// 		return c.Read(query)
// 	}

// 	var result = []m.Pattern{}
// 	for _, item := range models {
// 		var model, err = c.Read(&m.PatternQueryDto{ID: item.ID})
// 		if err != nil {
// 			return nil, err
// 		}

// 		result = append(result, model...)
// 	}
// 	return result, nil
// }

// func (c *PatternService) Delete(query *m.PatternQueryDto) (int, error) {
// 	var models []m.Pattern
// 	client, suffix := c.query(query, c.db)

// 	if err, _ := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil {
// 		return 0, err
// 	}

// 	for _, model := range models {
// 		c.recache(&model, c.keys(&model), true)
// 	}

// 	return len(models), client.Delete(&m.Pattern{}).Error
// }

package pods

import (
	"api/config"
	"api/helper"
	"api/logs"
	m "api/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type FullMetricsService struct {
	Pods    *PodsService
	Metrics *MetricsService
}

func NewFullMetricsService(db *gorm.DB, client *redis.Client, metrics *metrics.Clientset) *FullMetricsService {
	return &FullMetricsService{
		Pods:    NewPodsService(metrics),
		Metrics: NewMetricsService(db, client),
	}
}

type MetricsService struct {
	key string

	db     *gorm.DB
	client *redis.Client
}

func NewMetricsService(db *gorm.DB, client *redis.Client) *MetricsService {
	return &MetricsService{key: "METRICS", db: db, client: client}
}

func (c *MetricsService) precache(model *m.Metrics) {
	helper.Precache(c.client, c.key, fmt.Sprintf("ID=%d", model.ID), model)
	helper.Precache(c.client, c.key, fmt.Sprintf("PROJECT_ID=%d#NAME=%s#NAMESPACE=%s#CONTAINER_NAME=%s", model.ProjectID, model.Name, model.Namespace, model.ContainerName), model)

	var keys = []string{fmt.Sprintf("NAME=%s*", model.Name), fmt.Sprintf("NAMESPACE=%s*", model.Namespace), fmt.Sprintf("CONTAINER_NAME=%s*", model.ContainerName), fmt.Sprintf("PROJECT_ID=%d*", model.ProjectID), "", "CREATED_FROM=*", "CREATED_TO=*", "PAGE=*", "LIMIT=*"}
	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string, k string) interface{} {
			var data []m.Metrics
			if !strings.HasPrefix(str, "[") {
				data = make([]m.Metrics, 1)
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

				case "NAMESPACE":
					if model.Namespace != res[1] {
						return data
					}

				case "CONTAINER_NAME":
					if model.ContainerName != res[1] {
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

func (c *MetricsService) deepcache(models []m.Metrics, key string) interface{} {
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

	var items []m.Metrics
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

func (c *MetricsService) recache(model *m.Metrics, delete bool) {
	helper.Delcache(c.client, c.key, fmt.Sprintf("ID=%d*", model.ID))

	var keys = []string{fmt.Sprintf("NAME=%s*", model.Name), fmt.Sprintf("NAMESPACE=%s*", model.Namespace), fmt.Sprintf("CONTAINER_NAME=%s*", model.ContainerName), fmt.Sprintf("PROJECT_ID=%d*", model.ProjectID), "", "CREATED_FROM=*", "CREATED_TO=*", "PAGE=*", "LIMIT=*"}
	for _, key := range keys {
		helper.Recache(c.client, c.key, key, func(str string, suffix string) interface{} {
			if !strings.HasPrefix(str, "[") {
				str = fmt.Sprintf("[%s]", str)
			}

			var data []m.Metrics
			var result []m.Metrics

			json.Unmarshal([]byte(str), &data)
			for _, item := range data {
				if item.ID != model.ID {
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

func (c *MetricsService) query(dto *m.MetricsQueryDto, client *gorm.DB) (*gorm.DB, string) {
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

	if len(dto.Namespace) > 0 {
		suffix = append(suffix, fmt.Sprintf("NAMESPACE=%s", dto.Namespace))
		client = client.Where("namespace = ?", dto.Namespace)
	}

	if len(dto.ContainerName) > 0 {
		suffix = append(suffix, fmt.Sprintf("CONTAINER_NAME=%s", dto.ContainerName))
		client = client.Where("container_name = ?", dto.ContainerName)
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

func (c *MetricsService) Create(model *m.Metrics) error {
	// Check if such project_id exists
	if err, rows := helper.Getcache(c.db.Where("id = ?", model.ProjectID), c.client, "PROJECT", fmt.Sprintf("ID=%d", model.ProjectID), &[]m.Project{}); err != nil || rows == 0 {
		return fmt.Errorf("Requested project_id=%d do not exist", model.ProjectID)
	}

	var res *gorm.DB
	if res = c.db.Create(model); res.Error == nil {
		c.precache(model)
	}

	if res.Error != nil {
		go logs.DefaultLog("/controllers/metrics.go", res.Error)
		return fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	return nil
}

func (c *MetricsService) Read(query *m.MetricsQueryDto) ([]m.Metrics, error) {
	var model []m.Metrics
	client, suffix := c.query(query, c.db)

	err, _ := helper.Getcache(client.Order("created_at DESC"), c.client, c.key, suffix, &model)
	return model, err
}

func (c *MetricsService) Update(query *m.MetricsQueryDto, model *m.Metrics) ([]m.Metrics, error) {
	var res *gorm.DB
	client, suffix := c.query(query, c.db)

	var models = []m.Metrics{}
	if err, rows := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil || rows == 0 {
		return nil, fmt.Errorf("Requested model do not exist")
	}

	client, _ = c.query(query, c.db.Model(&m.Metrics{}))
	if res = client.Updates(model); res.Error != nil {
		go logs.DefaultLog("/controllers/metrics.go", res.Error)
		return nil, fmt.Errorf("Something unexpected happend: %v", res.Error)
	}

	for _, existed := range models {
		c.recache(existed.Fill(model), false)
	}

	return c.Read(query)
}

func (c *MetricsService) Delete(query *m.MetricsQueryDto) (int, error) {
	var models []m.Metrics
	client, suffix := c.query(query, c.db)

	if err, _ := helper.Getcache(client, c.client, c.key, suffix, &models); err != nil {
		return 0, err
	}

	for _, model := range models {
		c.recache(&model, true)
	}

	return len(models), client.Delete(&m.Metrics{}).Error
}

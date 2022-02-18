package pods

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/interfaces"
	"api/logs"
	"api/models"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type metricsController struct{}

func NewMetricsController() interfaces.Default {
	return &metricsController{}
}

func getScaledValue(q *resource.Quantity, scale int) (int64, int) {
	if scale >= 0 {
		return q.ScaledValue(resource.Scale(-int32(scale))), scale
	}

	var value int64
	if value = q.Value(); value%10 == 0 {
		return value, 0
	}

	for i := 1; i < 12; i++ {
		if value = q.ScaledValue(resource.Scale(-int32(i))); value%10 == 0 {
			return value, i
		}
	}

	return 0, 0
}

// @Tags Metrics
// @Summary Save an array of Pods Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Namespace of the Pod"
// @Param prefix query string false "Selector label, read more here: https://stackoverflow.com/a/47453572"
// @Param id path int true "Project primaray id"
// @Param namespace path string true "Namespace name"
// @Success 200 {object} models.Success{int}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics/{id}/{namespace}/ [post]
func (*metricsController) CreateAll(c *gin.Context) {
	var id int
	var namespace string

	if namespace = c.Param("namespace"); namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace shouldn't be empty")
		return
	}

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect project id param")
		return
	}

	options := metaV1.ListOptions{}
	if prefix := c.DefaultQuery("prefix", ""); prefix != "" {
		options.LabelSelector = fmt.Sprintf("app=%s", c.DefaultQuery("prefix", ""))
	}

	ctx := context.Background()
	result, err := config.Metrics.MetricsV1beta1().PodMetricses(namespace).List(ctx, options)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, pod := range result.Items {
		var count int
		key := fmt.Sprintf("METRICS:%s:%s:%d", pod.Namespace, pod.Name, id)
		if count, err = db.Redis.Get(ctx, key).Int(); err != nil {
			count = 0
		}

		if count < config.ENV.Metrics {
			db.Redis.Incr(ctx, key)

			for i, container := range pod.Containers {

				var cpuArg int64
				if cpuArg, err = db.Redis.Get(ctx, fmt.Sprintf("%s:CPU:%d", key, i)).Int64(); err != nil {
					cpuArg = 0
				}

				var cpuArgScale int
				if cpuArgScale, err = db.Redis.Get(ctx, fmt.Sprintf("%s:CPU:SCALE:%d", key, i)).Int(); err != nil {
					cpuArgScale = -1
				}

				var memoryArg int64
				if memoryArg, err = db.Redis.Get(ctx, fmt.Sprintf("%s:MEMORY:%d", key, i)).Int64(); err != nil {
					memoryArg = 0
				}

				var memoryArgScale int
				if memoryArgScale, err = db.Redis.Get(ctx, fmt.Sprintf("%s:MEMORY:SCALE:%d", key, i)).Int(); err != nil {
					memoryArgScale = -1
				}

				container.Usage.Cpu().MilliValue()

				cpu, cpuScale := getScaledValue(container.Usage.Cpu(), cpuArgScale)
				memory, memoryScale := getScaledValue(container.Usage.Memory(), memoryArgScale)

				db.Redis.Set(ctx, fmt.Sprintf("%s:CPU:%d", key, i), cpuArg+cpu/int64(config.ENV.Metrics), 0)
				db.Redis.Set(ctx, fmt.Sprintf("%s:MEMORY:%d", key, i), memoryArg+memory/int64(config.ENV.Metrics), 0)

				db.Redis.Set(ctx, fmt.Sprintf("%s:CPU:SCALE:%d", key, i), cpuScale, 0)
				db.Redis.Set(ctx, fmt.Sprintf("%s:MEMORY:SCALE:%d", key, i), memoryScale, 0)
			}

		} else {
			db.Redis.Del(ctx, key)
			model := make([]models.Metrics, len(pod.Containers))

			for i, container := range pod.Containers {
				model[i].ProjectID = uint32(id)

				model[i].Name = pod.Name
				model[i].Namespace = pod.Namespace
				model[i].ContainerName = container.Name

				model[i].CPU, _ = db.Redis.Get(ctx, fmt.Sprintf("%s:CPU:%d", key, i)).Int64()
				model[i].Memory, _ = db.Redis.Get(ctx, fmt.Sprintf("%s:MEMORY:%d", key, i)).Int64()

				cpuScale, _ := db.Redis.Get(ctx, fmt.Sprintf("%s:CPU:SCALE:%d", key, i)).Int()
				memScale, _ := db.Redis.Get(ctx, fmt.Sprintf("%s:MEMORY:SCALE:%d", key, i)).Int()

				model[i].CpuScale = int8(cpuScale)
				model[i].MemoryScale = int8(memScale)

				db.Redis.Del(ctx, fmt.Sprintf("%s:CPU:%d", key, i))
				db.Redis.Del(ctx, fmt.Sprintf("%s:MEMORY:%d", key, i))

				db.Redis.Del(ctx, fmt.Sprintf("%s:CPU:SCALE:%d", key, i))
				db.Redis.Del(ctx, fmt.Sprintf("%s:MEMORY:SCALE:%d", key, i))
			}

			if result := db.DB.Create(&model); result.Error != nil || result.RowsAffected == 0 {
				helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
				go logs.DefaultLog("/controllers/k3s/pods/metrics.go", result.Error)
				return
			}
		}
	}
}

// @Tags Metrics
// @Summary Save Pods Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Namespace of the Pod"
// @Param name path string true "Specified name of the Pod"
// @Param id path int true "Project primaray id"
// @Param namespace path string true "Namespace name"
// @Success 200 {object} models.Success{int}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics/{id}/{namespace}/{name} [post]
func (*metricsController) CreateOne(c *gin.Context) {
	var id int
	var name string
	var namespace string

	if name = c.Param("name"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	if namespace = c.Param("namespace"); namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace shouldn't be empty")
		return
	}

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect project id param")
		return
	}

	ctx := context.Background()
	result, err := config.Metrics.MetricsV1beta1().PodMetricses(namespace).Get(ctx, name, metaV1.GetOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	var count int
	key := fmt.Sprintf("METRICS:%s:%s:%d", result.Namespace, result.Name, id)
	if count, err = db.Redis.Get(ctx, key).Int(); err != nil {
		count = 0
	}

	if count < config.ENV.Metrics {
		db.Redis.Incr(ctx, key)

		for i, container := range result.Containers {

			var cpuArg int64
			if cpuArg, err = db.Redis.Get(ctx, fmt.Sprintf("%s:CPU:%d", key, i)).Int64(); err != nil {
				cpuArg = 0
			}

			var cpuArgScale int
			if cpuArgScale, err = db.Redis.Get(ctx, fmt.Sprintf("%s:CPU:SCALE:%d", key, i)).Int(); err != nil {
				cpuArgScale = -1
			}

			var memoryArg int64
			if memoryArg, err = db.Redis.Get(ctx, fmt.Sprintf("%s:MEMORY:%d", key, i)).Int64(); err != nil {
				memoryArg = 0
			}

			var memoryArgScale int
			if memoryArgScale, err = db.Redis.Get(ctx, fmt.Sprintf("%s:MEMORY:SCALE:%d", key, i)).Int(); err != nil {
				memoryArgScale = -1
			}

			container.Usage.Cpu().MilliValue()

			cpu, cpuScale := getScaledValue(container.Usage.Cpu(), cpuArgScale)
			memory, memoryScale := getScaledValue(container.Usage.Memory(), memoryArgScale)

			db.Redis.Set(ctx, fmt.Sprintf("%s:CPU:%d", key, i), cpuArg+cpu/int64(config.ENV.Metrics), 0)
			db.Redis.Set(ctx, fmt.Sprintf("%s:MEMORY:%d", key, i), memoryArg+memory/int64(config.ENV.Metrics), 0)

			db.Redis.Set(ctx, fmt.Sprintf("%s:CPU:SCALE:%d", key, i), cpuScale, 0)
			db.Redis.Set(ctx, fmt.Sprintf("%s:MEMORY:SCALE:%d", key, i), memoryScale, 0)
		}

	} else {
		db.Redis.Del(ctx, key)
		model := make([]models.Metrics, len(result.Containers))

		for i, container := range result.Containers {
			model[i].ProjectID = uint32(id)

			model[i].Name = result.Name
			model[i].Namespace = result.Namespace
			model[i].ContainerName = container.Name

			model[i].CPU, _ = db.Redis.Get(ctx, fmt.Sprintf("%s:CPU:%d", key, i)).Int64()
			model[i].Memory, _ = db.Redis.Get(ctx, fmt.Sprintf("%s:MEMORY:%d", key, i)).Int64()

			cpuScale, _ := db.Redis.Get(ctx, fmt.Sprintf("%s:CPU:SCALE:%d", key, i)).Int()
			memScale, _ := db.Redis.Get(ctx, fmt.Sprintf("%s:MEMORY:SCALE:%d", key, i)).Int()

			model[i].CpuScale = int8(cpuScale)
			model[i].MemoryScale = int8(memScale)

			db.Redis.Del(ctx, fmt.Sprintf("%s:CPU:%d", key, i))
			db.Redis.Del(ctx, fmt.Sprintf("%s:MEMORY:%d", key, i))

			db.Redis.Del(ctx, fmt.Sprintf("%s:CPU:SCALE:%d", key, i))
			db.Redis.Del(ctx, fmt.Sprintf("%s:MEMORY:SCALE:%d", key, i))
		}

		if result := db.DB.Create(&model); result.Error != nil || result.RowsAffected == 0 {
			helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
			go logs.DefaultLog("/controllers/k3s/pods/metrics.go", result.Error)
			return
		}
	}
}

// @Tags Metrics
// @Summary Get Pod Metrics by Project ID
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Project id"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} models.Success{result=[]models.Metrics}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics/{id} [get]
func (*metricsController) ReadOne(c *gin.Context) {
	var id int
	var model []models.Metrics

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}
	page, limit := helper.Pagination(c)

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("PROJECT_ID=%d", id)))
	if err := helper.PrecacheResult(fmt.Sprintf("METRICS:%s", hex.EncodeToString(hasher.Sum(nil))), db.DB.Where("project_id = ?", id).Order("created_at DESC").Offset(page*config.ENV.Items).Limit(limit), &model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/k3s/pods/metrics.go", err.Error())
		return
	}

	// TODO: Maybe one day ....
	// var items int64
	// var err error
	// if items, err = db.Redis.Get(context.Background(), "nLINK").Int64(); err != nil {
	// 	items = -1
	// 	go (&models.Link{}).Redis(db.DB, db.Redis)
	// 	go logs.DefaultLog("/controllers/k3s/pods/metrics.go", err.Error())
	// }

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: model,
		Items:  int64(len(model)),
		// TotalItems: items,
	})
}

func (*metricsController) ReadAll(c *gin.Context) {
}

func (*metricsController) UpdateOne(c *gin.Context) {}
func (*metricsController) UpdateAll(c *gin.Context) {}
func (*metricsController) DeleteOne(c *gin.Context) {}
func (*metricsController) DeleteAll(c *gin.Context) {}

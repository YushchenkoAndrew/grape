package pods

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/interfaces"
	"api/logs"
	"api/models"
	"context"
	"fmt"
	"net/http"
	"strconv"

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
	if value = q.Value(); value >= 100 {
		return value, 0
	}

	for i := 3; i < 12; i += 3 {
		if value = q.ScaledValue(resource.Scale(-int32(i))); value >= 100 {
			return value, i
		}
	}

	return 0, 0
}

func (*metricsController) CreateAll(c *gin.Context) {}

// @Tags Metrics
// @Summary Save Pods Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of Service"
// @Param namespace path string true "Namespace name"
// @Success 200 {object} models.Success{int}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics/{namespace}/{name} [get]
func (*metricsController) CreateOne(c *gin.Context) {
	var name string
	var namespace string

	if name = c.Param("name"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	if namespace = c.Param("namespace"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace shouldn't be empty")
		return
	}

	ctx := context.Background()
	result, err := config.Metrics.MetricsV1beta1().PodMetricses(namespace).Get(ctx, name, metaV1.GetOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	var count int
	key := fmt.Sprintf("METRICS:%s:%s", result.Namespace, result.Name)
	if count, err = db.Redis.Get(ctx, key).Int(); err != nil {
		count = 0
	}

	if count < config.ENV.Metrics {
		db.Redis.Incr(ctx, key)

		for index, container := range result.Containers {
			i := strconv.Itoa(index)

			var cpuArg int64
			if cpuArg, err = db.Redis.Get(ctx, key+":CPU:"+i).Int64(); err != nil {
				cpuArg = 0
			}

			var cpuArgScale int
			if cpuArgScale, err = db.Redis.Get(ctx, key+":CPU:SCALE:"+i).Int(); err != nil {
				cpuArgScale = -1
			}

			var memoryArg int64
			if memoryArg, err = db.Redis.Get(ctx, key+":MEMORY:"+i).Int64(); err != nil {
				memoryArg = 0
			}

			var memoryArgScale int
			if memoryArgScale, err = db.Redis.Get(ctx, key+":MEMORY:SCALE:"+i).Int(); err != nil {
				memoryArgScale = -1
			}

			container.Usage.Cpu().MilliValue()

			cpu, cpuScale := getScaledValue(container.Usage.Cpu(), cpuArgScale)
			memory, memoryScale := getScaledValue(container.Usage.Memory(), memoryArgScale)

			db.Redis.Set(ctx, key+":CPU:"+i, cpuArg+cpu/int64(config.ENV.Metrics), 0)
			db.Redis.Set(ctx, key+":MEMORY:"+i, memoryArg+memory/int64(config.ENV.Metrics), 0)

			db.Redis.Set(ctx, key+":CPU:SCALE:"+i, cpuScale, 0)
			db.Redis.Set(ctx, key+":MEMORY:SCALE:"+i, memoryScale, 0)
		}

	} else {
		db.Redis.Del(ctx, key)
		model := make([]models.Metrics, len(result.Containers))

		for i, container := range result.Containers {
			index := strconv.Itoa(i)

			model[i].Name = result.Name
			model[i].Namespace = result.Namespace
			model[i].ContainerName = container.Name

			model[i].CPU, _ = db.Redis.Get(ctx, key+":CPU:"+index).Int64()
			model[i].Memory, _ = db.Redis.Get(ctx, key+":MEMORY:"+index).Int64()

			cpuScale, _ := db.Redis.Get(ctx, key+":CPU:SCALE:"+index).Int()
			memScale, _ := db.Redis.Get(ctx, key+":MEMORY:SCALE:"+index).Int()

			model[i].CpuScale = uint8(cpuScale)
			model[i].MemoryScale = uint8(memScale)

			db.Redis.Del(ctx, key+":CPU:"+index)
			db.Redis.Del(ctx, key+":MEMORY:"+index)

			db.Redis.Del(ctx, key+":CPU:SCALE:"+index)
			db.Redis.Del(ctx, key+":MEMORY:SCALE:"+index)
		}

		if result := db.DB.Create(&model); result.Error != nil || result.RowsAffected == 0 {
			helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
			go logs.DefaultLog("/controllers/k3s/pods/metrics.go", result.Error)
			return
		}
	}
}

// @Tags Metrics
// @Summary Get Pod Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of Service"
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.PodMetrics}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics/{namespace}/{name} [get]
func (*metricsController) ReadOne(c *gin.Context) {

	// TODO: Think about this should I have this impl of Metrics or
	// simply request data from Database for spec pod

	// var name string
	// var namespace string

	// if name = c.Param("name"); name == "" {
	// 	helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
	// 	return
	// }

	// if namespace = c.Param("namespace"); name == "" {
	// 	helper.ErrHandler(c, http.StatusBadRequest, "Namespace shouldn't be empty")
	// 	return
	// }

	// ctx := context.Background()
	// result, err := config.Metrics.MetricsV1beta1().PodMetricses(namespace).Get(ctx, name, metaV1.GetOptions{})
	// if err != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// helper.ResHandler(c, http.StatusOK, models.Success{
	// 	Status: "OK",
	// 	Result: &[1]v1.PodMetrics{*result},
	// 	Items:  1,
	// })
}

// @Tags Metrics
// @Summary Get Pods Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace query string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.PodMetrics}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/pod/metrics [get]
func (*metricsController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.Metrics.MetricsV1beta1().PodMetricses(c.Param("namespace")).List(ctx, metaV1.ListOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: &result.Items,
		Items:  int64(len(result.Items)),
	})
}

func (*metricsController) UpdateOne(c *gin.Context) {}
func (*metricsController) UpdateAll(c *gin.Context) {}
func (*metricsController) DeleteOne(c *gin.Context) {}
func (*metricsController) DeleteAll(c *gin.Context) {}
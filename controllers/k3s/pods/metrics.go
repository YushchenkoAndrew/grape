package pods

import (
	"api/helper"
	"api/interfaces"
	m "api/models"
	"api/service/k3s/pods"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type metricsController struct {
	service *pods.FullMetricsService
}

func NewMetricsController(s *pods.FullMetricsService) interfaces.Default {
	return &metricsController{service: s}
}

// @Tags Metrics
// @Summary Save Pods Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int true "Project primaray id"
// @Param label path string true "LabelSelector, read more here: https://stackoverflow.com/a/47453572"
// @Param namespace path string true "Namespace name"
// @Success 200 {object} m.Success{result=[]m.Metrics}
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics/{namespace}/{name} [post]
func (o *metricsController) CreateOne(c *gin.Context) {
	var id = helper.GetID(c)
	var label = c.Param("label")
	var namespace = c.Param("namespace")

	if id == 0 || label == "" || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect params { id: %t, label: %t, namespace: %t }", id == 0, label == "", namespace == ""))
		return
	}

	metrics, err := o.service.Pods.Read(namespace, &m.K3sListQueryDto{LabelSelector: label})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, fmt.Sprintf("Hmmm k3s is broken, I guess: %v", err))
		return
	}

	var models = []m.Metrics{}
	for _, pod := range metrics.Items {
		for _, item := range pod.Containers {
			// NOTE: Maybe, maybe one day check on inf.dec value
			var cpu, _ = item.Usage.Cpu().AsInt64()
			var memory, _ = item.Usage.Memory().AsInt64()

			var model = m.Metrics{ProjectID: uint32(id), Name: pod.Name, Namespace: pod.Namespace, ContainerName: item.Name, CPU: cpu, Memory: memory}
			if err := o.service.Metrics.Create(&model); err != nil {
				helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
				return
			}

			models = append(models, model)
		}
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Metrics
// @Summary Save an array of Pods Metrics
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int true "Project primaray id"
// @Param namespace path string true "Namespace of the Pod"
// @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// @Param api_version query string false "APIVersion defines the versioned schema of this representation of an object."
// @Param label_selector query string false "A selector to restrict the list of returned objects by their labels."
// @Param field_selector query string false "A selector to restrict the list of returned objects by their fields."
// @Param watch query bool false "Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion."
// @Param allow_watch_bookmarks query bool false "allowWatchBookmarks requests watch events with type "BOOKMARK"."
// @Param resource_version query string false "resourceVersionMatch determines how resourceVersion is applied to list calls."
// @Param limit query int64 false "limit is a maximum number of responses to return for a list call."
// @Param continue query string false "The continue option should be set when retrieving more results from the server."
// @Success 200 {object} m.Success{result=[]m.Metrics}
// @Success 200 {object} m.Success{int}
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics/{namespace} [post]
func (o *metricsController) CreateAll(c *gin.Context) {
	var id = helper.GetID(c)
	var query = m.K3sListQueryDto{}
	var namespace = c.Param("namespace")

	if err := c.ShouldBindQuery(&query); err != nil || id == 0 || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect params { id: %t, query: %v, namespace: %t }", id == 0, err, namespace == ""))
		return
	}

	metrics, err := o.service.Pods.Read(namespace, &query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, fmt.Sprintf("Hmmm k3s is broken, I guess: %v", err))
		return
	}

	var models = []m.Metrics{}
	for _, pod := range metrics.Items {
		for _, item := range pod.Containers {
			// NOTE: Maybe, maybe one day check on inf.dec value
			var cpu, _ = item.Usage.Cpu().AsInt64()
			var memory, _ = item.Usage.Memory().AsInt64()

			var model = m.Metrics{ProjectID: uint32(id), Name: pod.Name, Namespace: pod.Namespace, ContainerName: item.Name, CPU: cpu, Memory: memory}
			if err := o.service.Metrics.Create(&model); err != nil {
				helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
				return
			}

			models = append(models, model)
		}
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
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
// @Success 200 {object} m.Success{result=[]m.Metrics}
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics/{id} [get]
func (*metricsController) ReadOne(c *gin.Context) {
	// var id int
	// var model []m.Metrics

	// if !helper.GetID(c, &id) {
	// 	helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
	// 	return
	// }
	// page, limit := helper.Pagination(c)

	// hasher := md5.New()
	// hasher.Write([]byte(fmt.Sprintf("PROJECT_ID=%d", id)))
	// if err := helper.PrecacheResult(fmt.Sprintf("METRICS:%s", hex.EncodeToString(hasher.Sum(nil))), db.DB.Where("project_id = ?", id).Order("created_at DESC").Offset(page*config.ENV.Items).Limit(limit), &model); err != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
	// 	go logs.DefaultLog("/controllers/k3s/pods/metrics.go", err.Error())
	// 	return
	// }

	// // TODO: Maybe one day ....
	// // var items int64
	// // var err error
	// // if items, err = db.Redis.Get(context.Background(), "nLINK").Int64(); err != nil {
	// // 	items = -1
	// // 	go (&m.Link{}).Redis(db.DB, db.Redis)
	// // 	go logs.DefaultLog("/controllers/k3s/pods/metrics.go", err.Error())
	// // }

	// helper.ResHandler(c, http.StatusOK, &m.Success{
	// 	Status: "OK",
	// 	Result: model,
	// 	Items:  int64(len(model)),
	// 	// TotalItems: items,
	// })
}

func (*metricsController) ReadAll(c *gin.Context) {
}

func (*metricsController) UpdateOne(c *gin.Context) {}
func (*metricsController) UpdateAll(c *gin.Context) {}
func (*metricsController) DeleteOne(c *gin.Context) {}
func (*metricsController) DeleteAll(c *gin.Context) {}

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
// @Router /k3s/pod/metrics/{namespace}/{label} [post]
func (o *metricsController) CreateOne(c *gin.Context) {
	var id = helper.GetID(c, "id")
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
	var id = helper.GetID(c, "id")
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
// @Summary Get Pod Metrics by ID
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Success 200 {object} m.Success{result=[]m.Metrics}
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics/{id} [get]
func (o *metricsController) ReadOne(c *gin.Context) {
	var id = helper.GetID(c, "id")

	if id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	models, err := o.service.Metrics.Read(&m.MetricsQueryDto{ID: uint32(id)})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Metrics
// @Summary Read Metrics by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param namespace query string false "Pod namespace"
// @Param container_name query string false "Container name"
// @Param project_id query string false "ProjectID: '1'"
// @Param created_from query string false "CreatedAt date >= start"
// @Param created_to query string false "CreatedAt date <= end"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} m.Success{result=[]m.Metrics}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics [get]
func (o *metricsController) ReadAll(c *gin.Context) {
	var query = m.MetricsQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Metrics.Read(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: models,
		Page:   query.Page,
		Limit:  query.Limit,
		Items:  len(models),
	})
}

// @Tags Metrics
// @Summary Update Metrics by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body m.MetricsDto true "Metrics Data"
// @Success 200 {object} m.Success{result=[]m.Metrics}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics/{id} [put]
func (o *metricsController) UpdateOne(c *gin.Context) {
	var body m.MetricsDto
	var id = helper.GetID(c, "id")

	if err := c.ShouldBind(&body); err != nil || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	models, err := o.service.Metrics.Update(&m.MetricsQueryDto{ID: uint32(id)}, &m.Metrics{CPU: body.CPU, Memory: body.Memory})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Metrics
// @Summary Update Metrics by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param namespace query string false "Pod namespace"
// @Param container_name query string false "Container name"
// @Param project_id query string false "ProjectID: '1'"
// @Param created_from query string false "CreatedAt date >= start"
// @Param created_to query string false "CreatedAt date <= end"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} m.Success{result=[]m.Metrics}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics [put]
func (o *metricsController) UpdateAll(c *gin.Context) {
	var query = m.MetricsQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	var body m.MetricsDto
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Metrics.Update(&query, &m.Metrics{CPU: body.CPU, Memory: body.Memory})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Metrics
// @Summary Delete Metrics by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics/{id} [delete]
func (o *metricsController) DeleteOne(c *gin.Context) {
	var id = helper.GetID(c, "id")

	if id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	items, err := o.service.Metrics.Delete(&m.MetricsQueryDto{ID: uint32(id)})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []string{},
		Items:  items,
	})
}

// @Tags Metrics
// @Summary Delete Metrics by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param namespace query string false "Pod namespace"
// @Param container_name query string false "Container name"
// @Param project_id query string false "ProjectID: '1'"
// @Param created_from query string false "CreatedAt date >= start"
// @Param created_to query string false "CreatedAt date <= end"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/pod/metrics [delete]
func (o *metricsController) DeleteAll(c *gin.Context) {
	var query = m.MetricsQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	items, err := o.service.Metrics.Delete(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []string{},
		Items:  items,
	})
}

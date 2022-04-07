package k3s

import (
	"api/helper"
	"api/interfaces"
	m "api/models"
	"api/service/k3s"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
)

type serviceController struct {
	service *k3s.ServiceService
}

func NewServiceController(s *k3s.ServiceService) interfaces.Default {
	return &serviceController{service: s}
}

// @Tags Service
// @Summary Create Service
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Namespace name"
// @Param model body v1.Service true "Service config json"
// @Success 201 {object} m.Success{result=[]v1.Service}
// @failure 400 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/service/{namespace} [post]
func (o *serviceController) CreateOne(c *gin.Context) {
	var body v1.Service
	var namespace = c.Param("namespace")

	if err := c.ShouldBind(&body); namespace == "" || err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect { namespace: %t, body: %t }", namespace == "", err != nil))
		return
	}

	result, err := o.service.Create(namespace, &body)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: []v1.Service{*result},
		Items:  1,
	})
}

// @Tags Service
// @Summary Create Service list
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Service Namespace"
// @Param model body []v1.Service true "Service config json"
// @Success 201 {object} m.Success{result=[]v1.Service}
// @failure 400 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/service/list/{namespace} [post]
func (o *serviceController) CreateAll(c *gin.Context) {
	var body []v1.Service
	var namespace = c.Param("namespace")

	if err := c.ShouldBind(&body); namespace == "" || err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect { namespace: %t, body: %v }", namespace == "", err))
		return
	}

	var models = []v1.Service{}
	for _, item := range body {
		result, err := o.service.Create(namespace, &item)
		if err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		models = append(models, *result)
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Service
// @Summary Get Service by :label
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param label path string true "Specified service label selector"
// @Param namespace path string false "Service Namespace"
// @Success 200 {object} m.Success{result=[]v1.Service}
// @failure 400 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/service/{namespace}/{label} [get]
func (o *serviceController) ReadOne(c *gin.Context) {
	var label = c.Param("label")
	var namespace = c.Param("namespace")

	if namespace == "" || label == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Param required params not presented: { namespace: %t, label: %t }", namespace == "", label == ""))
		return
	}

	result, err := o.service.Read(namespace, &m.K3sListQueryDto{LabelSelector: label})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: result.Items,
		Items:  len(result.Items),
	})
}

// @Tags Service
// @Summary Get Service by query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Service Namespace"
// @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// @Param api_version query string false "APIVersion defines the versioned schema of this representation of an object."
// @Param label_selector query string false "A selector to restrict the list of returned objects by their labels."
// @Param field_selector query string false "A selector to restrict the list of returned objects by their fields."
// @Param watch query bool false "Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion."
// @Param allow_watch_bookmarks query bool false "allowWatchBookmarks requests watch events with type "BOOKMARK"."
// @Param resource_version query string false "resourceVersionMatch determines how resourceVersion is applied to list calls."
// @Param limit query int64 false "limit is a maximum number of responses to return for a list call."
// @Param continue query string false "The continue option should be set when retrieving more results from the server."
// @Success 200 {object} m.Success{result=[]v1.Service}
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/service/{namespace} [get]
func (o *serviceController) ReadAll(c *gin.Context) {
	var query = m.K3sListQueryDto{}
	var namespace = c.Param("namespace")

	if err := c.ShouldBindQuery(&query); err != nil || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	result, err := o.service.Read(namespace, &query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: result.Items,
		Items:  len(result.Items),
	})
}

// @Tags Service
// @Summary Update Service
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Service Namespace"
// @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// @Param api_version query string false "APIVersion defines the versioned schema of this representation of an object."
// @Param field_manager query string false "fieldManager is a name associated with the actor or entity that is making these changes"
// @Param field_validation query string false "fieldValidation determines how the server should respond to unknown/duplicate fields in the object in the request."
// @Param model body v1.Service true "Updated Service body. Field Name should be presented !!!"
// @Success 200 {object} m.Success{result=[]v1.Service}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/service/{namespace} [put]
func (o *serviceController) UpdateOne(c *gin.Context) {
	var query = m.K3sUpdateQueryDto{}
	var namespace = c.Param("namespace")

	if err := c.ShouldBindQuery(&query); err != nil || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad query request {query: %v, namespace: %t}", err, namespace == ""))
		return
	}

	var body v1.Service
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect body: %v", err))
		return
	}

	result, err := o.service.Update(namespace, &query, &body)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []v1.Service{*result},
		Items:  1,
	})
}

// @Tags Service
// @Summary Update Service list
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Service Namespace"
// @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// @Param api_version query string false "APIVersion defines the versioned schema of this representation of an object."
// @Param field_manager query string false "fieldManager is a name associated with the actor or entity that is making these changes"
// @Param field_validation query string false "fieldValidation determines how the server should respond to unknown/duplicate fields in the object in the request."
// @Param model body []v1.Service true "Updated Service body. Field Name should be presented !!!"
// @Success 200 {object} m.Success{result=[]v1.Service}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/service/list/{namespace} [put]
func (o *serviceController) UpdateAll(c *gin.Context) {
	var query = m.K3sUpdateQueryDto{}
	var namespace = c.Param("namespace")

	if err := c.ShouldBindQuery(&query); err != nil || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad query request {query: %v, namespace: %t}", err, namespace == ""))
		return
	}

	var body []v1.Service
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect body: %v", err))
		return
	}

	var models = []v1.Service{}
	for _, item := range body {
		result, err := o.service.Update(namespace, &query, &item)
		if err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		models = append(models, *result)
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Service
// @Summary Delete Service by :name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name query string true "Service name"
// @Param namespace path string false "Service Namespace"
// @Success 200 {object} m.Success{result=[]string}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/service/{namespace}/{name} [delete]
func (o *serviceController) DeleteOne(c *gin.Context) {
	var name = c.Param("name")
	var namespace = c.Param("namespace")

	if name == "" || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect params: { name: %t, namespace: %t }", name == "", namespace == ""))
		return
	}

	if err := o.service.Delete(namespace, name); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []string{},
		Items:  1,
	})
}

// @Tags Service
// @Summary Delete Service list
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Service Namespace"
// @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// @Param api_version query string false "APIVersion defines the versioned schema of this representation of an object."
// @Param label_selector query string false "A selector to restrict the list of returned objects by their labels."
// @Param field_selector query string false "A selector to restrict the list of returned objects by their fields."
// @Param watch query bool false "Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion."
// @Param allow_watch_bookmarks query bool false "allowWatchBookmarks requests watch events with type "BOOKMARK"."
// @Param resource_version query string false "resourceVersionMatch determines how resourceVersion is applied to list calls."
// @Param limit query int64 false "limit is a maximum number of responses to return for a list call."
// @Param continue query string false "The continue option should be set when retrieving more results from the server."
// @Success 200 {object} m.Success{result=[]string}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/service/{namespace} [delete]
func (o *serviceController) DeleteAll(c *gin.Context) {
	var query = m.K3sListQueryDto{}
	var namespace = c.Param("namespace")

	if err := c.ShouldBindQuery(&query); err != nil || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Read(namespace, &query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, item := range models.Items {
		if err := o.service.Delete(namespace, item.GetName()); err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []string{},
		Items:  len(models.Items),
	})
}

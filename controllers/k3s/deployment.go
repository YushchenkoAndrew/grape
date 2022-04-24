package k3s

import (
	"api/helper"
	"api/interfaces/controller"
	m "api/models"
	"api/service/k3s"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
)

type deploymentController struct {
	service *k3s.DeploymentService
}

func NewDeploymentController(s *k3s.DeploymentService) controller.Default {
	return &deploymentController{service: s}
}

// @Tags Deployment
// @Summary Create Deployment
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Deployment Namespace"
// @Param model body v1.Deployment true "Deployment config json"
// @Success 201 {object} m.Success{result=[]v1.Deployment}
// @failure 400 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/{namespace} [post]
func (o *deploymentController) CreateOne(c *gin.Context) {
	var body v1.Deployment
	var namespace = c.Param("namespace")

	if err := c.ShouldBind(&body); namespace == "" || err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect { namespace: %t, body: %v }", namespace == "", err))
		return
	}

	result, err := o.service.Create(namespace, &body)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: []v1.Deployment{*result},
		Items:  1,
	})
}

// @Tags Deployment
// @Summary Create Deployment list
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Deployment Namespace"
// @Param model body []v1.Deployment true "Deployment config json"
// @Success 201 {object} m.Success{result=[]v1.Deployment}
// @failure 400 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/list/{namespace} [post]
func (o *deploymentController) CreateAll(c *gin.Context) {
	var body []v1.Deployment
	var namespace = c.Param("namespace")

	if err := c.ShouldBind(&body); namespace == "" || err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect { namespace: %t, body: %v }", namespace == "", err))
		return
	}

	var models = []v1.Deployment{}
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

// @Tags Deployment
// @Summary Get Deployments by :label
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param label path string true "Specified deployment label selector"
// @Param namespace path string false "Deployment Namespace"
// @Success 200 {object} m.Success{result=[]v1.Deployment}
// @failure 400 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/{namespace}/{label} [get]
func (o *deploymentController) ReadOne(c *gin.Context) {
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

// @Tags Deployment
// @Summary Get Deployments by query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Deployment Namespace"
// @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// @Param api_version query string false "APIVersion defines the versioned schema of this representation of an object."
// @Param label_selector query string false "A selector to restrict the list of returned objects by their labels."
// @Param field_selector query string false "A selector to restrict the list of returned objects by their fields."
// @Param watch query bool false "Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion."
// @Param allow_watch_bookmarks query bool false "allowWatchBookmarks requests watch events with type "BOOKMARK"."
// @Param resource_version query string false "resourceVersionMatch determines how resourceVersion is applied to list calls."
// @Param limit query int64 false "limit is a maximum number of responses to return for a list call."
// @Param continue query string false "The continue option should be set when retrieving more results from the server."
// @Success 200 {object} m.Success{result=[]v1.Deployment}
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/{namespace} [get]
func (o *deploymentController) ReadAll(c *gin.Context) {
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

// @Tags Deployment
// @Summary Update Deployments
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Deployment Namespace"
// @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// @Param api_version query string false "APIVersion defines the versioned schema of this representation of an object."
// @Param field_manager query string false "fieldManager is a name associated with the actor or entity that is making these changes"
// @Param field_validation query string false "fieldValidation determines how the server should respond to unknown/duplicate fields in the object in the request."
// @Param model body v1.Deployment true "Updated Deployment body. Field Name should be presented !!!"
// @Success 200 {object} m.Success{result=[]v1.Deployment}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/{namespace} [put]
func (o *deploymentController) UpdateOne(c *gin.Context) {
	var query = m.K3sUpdateQueryDto{}
	var namespace = c.Param("namespace")

	if err := c.ShouldBindQuery(&query); err != nil || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad query request {query: %v, namespace: %t}", err, namespace == ""))
		return
	}

	var body v1.Deployment
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
		Result: []v1.Deployment{*result},
		Items:  1,
	})
}

// @Tags Deployment
// @Summary Update Deployments list
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Deployment Namespace"
// @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// @Param api_version query string false "APIVersion defines the versioned schema of this representation of an object."
// @Param field_manager query string false "fieldManager is a name associated with the actor or entity that is making these changes"
// @Param field_validation query string false "fieldValidation determines how the server should respond to unknown/duplicate fields in the object in the request."
// @Param model body []v1.Deployment true "Updated Deployment body. Field Name should be presented !!!"
// @Success 200 {object} m.Success{result=[]v1.Deployment}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/list/{namespace} [put]
func (o *deploymentController) UpdateAll(c *gin.Context) {
	var query = m.K3sUpdateQueryDto{}
	var namespace = c.Param("namespace")

	if err := c.ShouldBindQuery(&query); err != nil || namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad query request {query: %v, namespace: %t}", err, namespace == ""))
		return
	}

	var body []v1.Deployment
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect body: %v", err))
		return
	}

	var models = []v1.Deployment{}
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

// @Tags Deployment
// @Summary Delete Deployments by :name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name query string true "Deployment name"
// @Param namespace path string false "Deployment Namespace"
// @Success 200 {object} m.Success{result=[]string}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/{namespace}/{name} [delete]
func (o *deploymentController) DeleteOne(c *gin.Context) {
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

// @Tags Deployment
// @Summary Delete Deployments list
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Deployment Namespace"
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
// @Router /k3s/deployment/{namespace} [delete]
func (o *deploymentController) DeleteAll(c *gin.Context) {
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

package namespace

// import (
// 	"grape/helper"
// 	"grape/interfaces/controller"
// 	m "grape/models"
// 	"grape/service/k3s"
// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	v1 "k8s.io/api/core/v1"
// )

// type namespaceController struct {
// 	service *k3s.NamespaceService
// }

// func NewNamespaceController(s *k3s.NamespaceService) controller.Default {
// 	return &namespaceController{service: s}
// }

// // @Tags Namespace
// // @Summary Create Namespace
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param model body v1.Namespace true "Namespace json body"
// // @Success 201 {object} m.Success{result=[]v1.Namespace}
// // @failure 400 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /k3s/namespace [post]
// func (o *namespaceController) CreateOne(c *gin.Context) {
// 	var body v1.Namespace
// 	if err := c.ShouldBind(&body); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect body is not setted: %v", err))
// 		return
// 	}

// 	result, err := o.service.Create(&body)
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusCreated, &m.Success{
// 		Status: "OK",
// 		Result: []v1.Namespace{*result},
// 		Items:  1,
// 	})
// }

// // @Tags Namespace
// // @Summary Create list of Namespace
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param model body []v1.Namespace true "Namespace json List"
// // @Success 201 {object} m.Success{result=[]v1.Namespace}
// // @failure 400 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /k3s/namespace/list [post]
// func (o *namespaceController) CreateAll(c *gin.Context) {
// 	var body []v1.Namespace
// 	if err := c.ShouldBind(&body); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect body is not setted: %v", err))
// 		return
// 	}

// 	var models = []v1.Namespace{}
// 	for _, item := range body {
// 		result, err := o.service.Create(&item)
// 		if err != nil {
// 			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		models = append(models, *result)
// 	}

// 	helper.ResHandler(c, http.StatusCreated, &m.Success{
// 		Status: "OK",
// 		Result: models,
// 		Items:  len(models),
// 	})
// }

// // @Tags Namespace
// // @Summary Get Namespace by :label
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param label path string true "A selector to restrict the list of returned objects by their labels."
// // @Success 200 {object} m.Success{result=[]v1.Namespace}
// // @failure 400 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /k3s/namespace/{label} [get]
// func (o *namespaceController) ReadOne(c *gin.Context) {
// 	var label = c.Param("label")
// 	if label == "" {
// 		helper.ErrHandler(c, http.StatusBadRequest, "Label shouldn't be empty")
// 		return
// 	}

// 	result, err := o.service.Read(&m.K3sListQueryDto{LabelSelector: label})
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: result.Items,
// 		Items:  len(result.Items),
// 	})
// }

// // @Tags Namespace
// // @Summary Get Namespace by query
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// // @Param api_version query string false "grapeVersion defines the versioned schema of this representation of an object."
// // @Param label_selector query string false "A selector to restrict the list of returned objects by their labels."
// // @Param field_selector query string false "A selector to restrict the list of returned objects by their fields."
// // @Param watch query bool false "Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion."
// // @Param allow_watch_bookmarks query bool false "allowWatchBookmarks requests watch events with type "BOOKMARK"."
// // @Param resource_version query string false "resourceVersionMatch determines how resourceVersion is applied to list calls."
// // @Param limit query int64 false "limit is a maximum number of responses to return for a list call."
// // @Param continue query string false "The continue option should be set when retrieving more results from the server."
// // @Success 200 {object} m.Success{result=[]v1.Namespace}
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /k3s/namespace [get]
// func (o *namespaceController) ReadAll(c *gin.Context) {
// 	var query = m.K3sListQueryDto{}

// 	if err := c.ShouldBindQuery(&query); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
// 		return
// 	}

// 	result, err := o.service.Read(&query)
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: result.Items,
// 		Items:  len(result.Items),
// 	})
// }

// // @Tags Namespace
// // @Summary Update Namespace
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// // @Param api_version query string false "grapeVersion defines the versioned schema of this representation of an object."
// // @Param field_manager query string false "fieldManager is a name associated with the actor or entity that is making these changes"
// // @Param field_validation query string false "fieldValidation determines how the server should respond to unknown/duplicate fields in the object in the request."
// // @Param model body v1.Namespace true "Updated Namespace body. Field Name should be presented !!!"
// // @Success 200 {object} m.Success{result=[]v1.Namespace}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /k3s/namespace [put]
// func (o *namespaceController) UpdateOne(c *gin.Context) {
// 	var query = m.K3sUpdateQueryDto{}
// 	if err := c.ShouldBindQuery(&query); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad query request: %v", err))
// 		return
// 	}

// 	var body v1.Namespace
// 	if err := c.ShouldBind(&body); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect body: %v", err))
// 		return
// 	}

// 	result, err := o.service.Update(&query, &body)
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: []v1.Namespace{*result},
// 		Items:  1,
// 	})
// }

// // @Tags Namespace
// // @Summary Update Namespace list
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// // @Param api_version query string false "grapeVersion defines the versioned schema of this representation of an object."
// // @Param field_manager query string false "fieldManager is a name associated with the actor or entity that is making these changes"
// // @Param field_validation query string false "fieldValidation determines how the server should respond to unknown/duplicate fields in the object in the request."
// // @Param model body []v1.Namespace true "Updated Namespace body. Field Name should be presented !!!"
// // @Success 200 {object} m.Success{result=[]v1.Namespace}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /k3s/namespace/list [put]
// func (o *namespaceController) UpdateAll(c *gin.Context) {
// 	var query = m.K3sUpdateQueryDto{}
// 	if err := c.ShouldBindQuery(&query); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad query request: %v", err))
// 		return
// 	}

// 	var body []v1.Namespace
// 	if err := c.ShouldBind(&body); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Incorrect body: %v", err))
// 		return
// 	}

// 	var models = []v1.Namespace{}
// 	for _, item := range body {
// 		result, err := o.service.Update(&query, &item)
// 		if err != nil {
// 			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		models = append(models, *result)
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: models,
// 		Items:  len(models),
// 	})
// }

// // @Tags Namespace
// // @Summary Delete Namespace by :name
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param name query string true "Namespace name"
// // @Success 200 {object} m.Success{result=[]string}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /k3s/namespace/{name} [delete]
// func (o *namespaceController) DeleteOne(c *gin.Context) {
// 	var name = c.Param("name")
// 	if name == "" {
// 		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
// 		return
// 	}

// 	if err := o.service.Delete(name); err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: []string{},
// 		Items:  1,
// 	})
// }

// // @Tags Namespace
// // @Summary Delete Namespace by query
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param kind query string false "Kind is a string value representing the REST resource this object represents."
// // @Param api_version query string false "grapeVersion defines the versioned schema of this representation of an object."
// // @Param label_selector query string false "A selector to restrict the list of returned objects by their labels."
// // @Param field_selector query string false "A selector to restrict the list of returned objects by their fields."
// // @Param watch query bool false "Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion."
// // @Param allow_watch_bookmarks query bool false "allowWatchBookmarks requests watch events with type "BOOKMARK"."
// // @Param resource_version query string false "resourceVersionMatch determines how resourceVersion is applied to list calls."
// // @Param limit query int64 false "limit is a maximum number of responses to return for a list call."
// // @Param continue query string false "The continue option should be set when retrieving more results from the server."
// // @Success 200 {object} m.Success{result=[]string}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /k3s/namespace [delete]
// func (o *namespaceController) DeleteAll(c *gin.Context) {
// 	var query = m.K3sListQueryDto{}

// 	if err := c.ShouldBindQuery(&query); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
// 		return
// 	}

// 	models, err := o.service.Read(&query)
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	for _, item := range models.Items {
// 		if err := o.service.Delete(item.GetName()); err != nil {
// 			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: []string{},
// 		Items:  len(models.Items),
// 	})
// }

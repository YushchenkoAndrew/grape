package k3s

import (
	"api/helper"
	"api/interfaces"
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

func NewDeploymentController(s *k3s.DeploymentService) interfaces.Default {
	return &deploymentController{service: s}
}

func (*deploymentController) CreateAll(c *gin.Context) {}

// @Tags Deployment
// @Summary Create Deployment
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Namespace name"
// @Param model body v1.Deployment true "Deployment config file"
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
		Result: []v1.Deployment{*result},
		Items:  1,
	})
}

// @Tags Deployment
// @Summary Get Deployments
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of deployment"
// @Param namespace path string false "Namespace name"
// @Success 200 {object} m.Success{result=[]v1.Deployment}
// @failure 400 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/{namespace}/{name} [get]
func (o *deploymentController) ReadOne(c *gin.Context) {
	var namespace = c.Param("namespace")
	if namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Param namespace is required"))
		return
	}

	result, err := o.service.Read(namespace, &m.K3sQueryDto{LabelSelector: c.Param("name")})
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
// @Summary Get Deployments
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Namespace name"
// @Success 200 {object} m.Success{result=[]v1.Deployment}
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /k3s/deployment/{namespace} [get]
func (o *deploymentController) ReadAll(c *gin.Context) {
	var query = m.K3sQueryDto{}
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

func (*deploymentController) UpdateAll(c *gin.Context) {}
func (*deploymentController) UpdateOne(c *gin.Context) {}
func (*deploymentController) DeleteAll(c *gin.Context) {}
func (*deploymentController) DeleteOne(c *gin.Context) {}

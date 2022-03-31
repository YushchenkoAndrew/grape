package k3s

import (
	"api/config"
	"api/helper"
	"api/interfaces"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type serviceController struct{}

func NewServiceController() interfaces.Default {
	return &serviceController{}
}

func (*serviceController) CreateAll(c *gin.Context) {}

// @Tags Service
// @Summary Create Service
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string true "Namespace name"
// @Param model body v1.Service true "Deployment config file"
// @Success 201 {object} models.Success{result=[]v1.Service}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/service/{namespace} [post]
func (*serviceController) CreateOne(c *gin.Context) {
	var namespace = c.Param("namespace")
	if namespace == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Namespace name shouldn't be empty")
		return
	}

	var body v1.Service
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body is not setted")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.CoreV1().Services(namespace).Create(ctx, &body, metaV1.CreateOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &models.Success{
		Status: "OK",
		Result: []v1.Service{*result},
		Items:  1,
	})
}

// @Tags Service
// @Summary Get Service
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Specified name of Service"
// @Param namespace path string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Service}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/service/{namespace}/{name} [get]
func (*serviceController) ReadOne(c *gin.Context) {
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

	ctx := context.Background()
	result, err := config.K3s.CoreV1().Services(namespace).Get(ctx, name, metaV1.GetOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: []v1.Service{*result},
		Items:  1,
	})
}

// @Tags Service
// @Summary Get Service
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param namespace path string false "Namespace name"
// @Success 200 {object} models.Success{result=[]v1.Service}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/service/{namespace} [get]
func (*serviceController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.K3s.CoreV1().Services(c.Param("namespace")).List(ctx, metaV1.ListOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: result.Items,
		Items:  len(result.Items),
	})
}

func (*serviceController) UpdateAll(c *gin.Context) {}
func (*serviceController) UpdateOne(c *gin.Context) {}
func (*serviceController) DeleteAll(c *gin.Context) {}
func (*serviceController) DeleteOne(c *gin.Context) {}

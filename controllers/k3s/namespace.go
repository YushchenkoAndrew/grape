package k3s

import (
	"api/config"
	"api/helper"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceController struct{}

// @Tags K3s
// @Summary Create Namespace
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param model body v1.Namespace true "Link info"
// @Success 201 {object} models.Success{result=[]v1.Namespace}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/namespace [post]
func (*NamespaceController) Create(c *gin.Context) {
	var body v1.Namespace
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body is not setted")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.CoreV1().Namespaces().Create(ctx, &body, metaV1.CreateOptions{})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status: "OK",
		Result: &[1]v1.Namespace{*result},
		Items:  1,
	})
}

// @Tags K3s
// @Summary Get Deployments
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name query string false "Specified name of Namespace"
// @Success 200 {object} models.Success{result=[]v1.Namespace}
// @failure 400 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/namespace/{name} [get]
func (*NamespaceController) ReadOne(c *gin.Context) {
	var name = c.Param("name")
	if name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	ctx := context.Background()
	result, err := config.K3s.CoreV1().Namespaces().Get(ctx, name, metaV1.GetOptions{})

	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status: "OK",
		Result: &[1]v1.Namespace{*result},
		Items:  1,
	})
}

// @Tags K3s
// @Summary Get Deployments
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Success 200 {object} models.Success{result=[]v1.Namespace}
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /k3s/namespace [get]
func (*NamespaceController) ReadAll(c *gin.Context) {
	ctx := context.Background()
	result, err := config.K3s.CoreV1().Namespaces().List(ctx, metaV1.ListOptions{})

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
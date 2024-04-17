package src

import (
	"grape/src/common/dto/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexT interface {
	Ping(c *gin.Context)
}

type indexController struct {
	service *indexService
}

func NewIndexController(s *indexService) IndexT {
	return &indexController{service: s}
}

// @Summary Ping/Pong
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Success 200 {object} interface{}
// @Router /ping [get]
func (*indexController) Ping(ctx *gin.Context) {
	response.Handler(ctx, http.StatusOK, gin.H{"message": "pong"}, nil)
}

package src

import (
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
// @Success 200 {object} interface{}
// @Router /ping [get]
func (*indexController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

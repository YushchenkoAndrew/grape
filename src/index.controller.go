package src

import (
	m "grape/models"
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
// @Success 200 {object} m.Ping
// @failure 429 {object} m.Error
// @Router /ping [get]
func (*indexController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, m.Ping{
		Status:  "OK",
		Message: "pong",
	})
}

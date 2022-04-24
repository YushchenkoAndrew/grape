package controllers

import (
	"api/helper"
	"api/interfaces/controller"
	m "api/models"
	"api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type botController struct {
	service *service.BotService
}

func NewBotController(s *service.BotService) controller.Bot {
	return &botController{service: s}
}

// @Tags Bot
// @Summary Execute redis Command from request
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body m.BotRedis true "Redis Command"
// @Success 200 {object} m.DefultRes
// @failure 400 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /bot/redis [post]
func (o *botController) Redis(c *gin.Context) {
	var body m.BotRedisDto

	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body")
		return
	}

	res, err := o.service.RunRedis(&body)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, m.DefultRes{
		Status:  "OK",
		Message: "Success",
		Result:  res,
	})
}

package routes

import (
	c "api/controllers"
	"api/interfaces"
	m "api/middleware"

	"github.com/gin-gonic/gin"
)

type botRouter struct {
	auth *gin.RouterGroup
	bot  interfaces.Bot
}

func NewBotRouter(rg *gin.RouterGroup) interfaces.Router {
	return &botRouter{
		auth: rg.Group("/bot", m.GetMiddleware().Auth()),
		bot:  c.NewBotController(),
	}
}

func (c *botRouter) Init() {
	c.auth.POST("/redis", c.bot.Redis)
}

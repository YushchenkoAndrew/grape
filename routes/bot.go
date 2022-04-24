package routes

import (
	c "api/controllers"
	"api/interfaces"
	i "api/interfaces/controller"
	m "api/middleware"
	s "api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type botRouter struct {
	auth *gin.RouterGroup
	bot  i.Bot
}

func NewBotRouter(rg *gin.RouterGroup, db *gorm.DB, client *redis.Client) interfaces.Router {
	return &botRouter{
		auth: rg.Group("/bot", m.GetMiddleware().Auth()),
		bot:  c.NewBotController(s.NewBotService(db, client)),
	}
}

func (c *botRouter) Init() {
	c.auth.POST("/redis", c.bot.Redis)
}

package routes

import (
	c "api/controllers"
	"api/interfaces"
	m "api/middleware"
	s "api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type linkRouter struct {
	route *gin.RouterGroup
	auth  *gin.RouterGroup
	link  interfaces.Default
}

func NewLinkRouter(rg *gin.RouterGroup, db *gorm.DB, client *redis.Client) interfaces.Router {
	return &linkRouter{
		route: rg.Group(("/link")),
		auth:  rg.Group("/link", m.GetMiddleware().Auth()),
		link:  c.NewLinkController(s.NewLinkService(db, client)),
	}
}

func (c *linkRouter) Init() {
	c.auth.POST("/list/:id", c.link.CreateAll)
	c.auth.POST("/:id", c.link.CreateOne)

	c.route.GET("/:id", c.link.ReadOne)
	c.route.GET("", c.link.ReadAll)

	c.auth.PUT("/:id", c.link.UpdateOne)
	c.auth.PUT("", c.link.UpdateAll)

	c.auth.DELETE("/:id", c.link.DeleteOne)
	c.auth.DELETE("", c.link.DeleteAll)
}

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

type patternRouter struct {
	route   *gin.RouterGroup
	auth    *gin.RouterGroup
	pattern i.Default
}

func NewPatternRouter(rg *gin.RouterGroup, db *gorm.DB, client *redis.Client) interfaces.Router {
	return &patternRouter{
		route:   rg.Group(("/pattern")),
		auth:    rg.Group("/pattern", m.GetMiddleware().Auth()),
		pattern: c.NewPatternController(s.NewPatternService(db, client)),
	}
}

func (c *patternRouter) Init() {
	c.auth.POST("/list/:id", c.pattern.CreateAll)
	c.auth.POST("/:id", c.pattern.CreateOne)

	c.route.GET("/:id", c.pattern.ReadOne)
	c.route.GET("", c.pattern.ReadAll)

	c.auth.PUT("/:id", c.pattern.UpdateOne)
	c.auth.PUT("", c.pattern.UpdateAll)

	c.auth.DELETE("/:id", c.pattern.DeleteOne)
	c.auth.DELETE("", c.pattern.DeleteAll)
}

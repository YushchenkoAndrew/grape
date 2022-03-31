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

type projectRouter struct {
	route   *gin.RouterGroup
	auth    *gin.RouterGroup
	project interfaces.Default
}

func NewProjectRouter(rg *gin.RouterGroup, db *gorm.DB, client *redis.Client) interfaces.Router {
	return &projectRouter{
		route:   rg.Group(("/project")),
		auth:    rg.Group("/project", m.GetMiddleware().Auth()),
		project: c.NewProjectController(s.NewFullProjectService(db, client)),
	}
}

func (c *projectRouter) Init() {
	c.auth.POST("", c.project.CreateOne)
	c.auth.POST("/list", c.project.CreateAll)

	c.route.GET("/:name", c.project.ReadOne)
	c.route.GET("", c.project.ReadAll)

	c.auth.PUT("/:name", c.project.UpdateOne)
	c.auth.PUT("", c.project.UpdateAll)

	c.auth.DELETE("/:name", c.project.DeleteOne)
	c.auth.DELETE("", c.project.DeleteAll)
}

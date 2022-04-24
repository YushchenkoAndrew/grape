package routes

import (
	"api/interfaces"
	m "api/middleware"

	"github.com/gin-gonic/gin"
)

type worldRouter struct {
	route *gin.RouterGroup
	auth  *gin.RouterGroup
	// world interfaces.Default
}

func NewWorldRouter(rg *gin.RouterGroup) interfaces.Router {
	return &worldRouter{
		route: rg.Group(("/world")),
		auth:  rg.Group("/world", m.GetMiddleware().Auth()),
		// world: c.NewWorldController(),
	}
}

// FIXME: !!
func (c *worldRouter) Init() {
	// c.auth.POST("", c.world.CreateOne)
	// c.auth.POST("/list", c.world.CreateAll)

	// c.route.GET("/:id", c.world.ReadOne)
	// c.route.GET("", c.world.ReadAll)

	// c.auth.PUT("/:id", c.world.UpdateOne)
	// c.auth.PUT("", c.world.UpdateAll)

	// c.auth.DELETE("/:id", c.world.DeleteOne)
	// c.auth.DELETE("", c.world.DeleteAll)
}

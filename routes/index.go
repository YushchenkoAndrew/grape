package routes

import (
	c "api/controllers"
	"api/interfaces"
	s "api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type indexRouter struct {
	route     *gin.RouterGroup
	index     interfaces.Index
	subRoutes *[]interfaces.Router
}

func NewIndexRouter(route *gin.RouterGroup, subRoutes *[]interfaces.Router, db *gorm.DB, client *redis.Client) interfaces.Router {
	return &indexRouter{
		route:     route,
		index:     c.NewIndexController(s.NewIndexService(db, client)),
		subRoutes: subRoutes,
	}
}

func (c *indexRouter) Init() {
	c.route.GET("/ping", c.index.Ping)
	c.route.GET("/trace/:ip", c.index.TraceIp)
	c.route.POST("/login", c.index.Login)
	c.route.POST("/refresh", c.index.Refresh)

	for _, route := range *c.subRoutes {
		route.Init()
	}
}

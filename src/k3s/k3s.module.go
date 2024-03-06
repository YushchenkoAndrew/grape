package k3s

// import (
// 	"grape/interfaces"
// 	m "grape/middleware"

// 	"github.com/gin-gonic/gin"
// )

// type k3sRouter struct {
// 	auth      *gin.RouterGroup
// 	subRoutes []interfaces.Router
// }

// func NewK3sRouter(rg *gin.RouterGroup, handlers []func(*gin.RouterGroup) interfaces.Router) interfaces.Router {
// 	route := rg.Group("/k3s")
// 	var subRoutes []interfaces.Router
// 	for _, handler := range handlers {
// 		subRoutes = append(subRoutes, handler(route))
// 	}

// 	return &k3sRouter{
// 		auth:      rg.Group("/k3s", m.GetMiddleware().Auth()),
// 		subRoutes: subRoutes,
// 	}
// }

// func (c *k3sRouter) Init() {
// 	for _, route := range c.subRoutes {
// 		route.Init()
// 	}
// }

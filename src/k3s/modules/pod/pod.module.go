package pod

// import (
// 	c "grape/controllers/k3s"
// 	"grape/interfaces"
// 	"grape/interfaces/k3s"
// 	m "grape/middleware"

// 	"github.com/gin-gonic/gin"
// )

// type podsRouter struct {
// 	auth      *gin.RouterGroup
// 	pods      k3s.Pods
// 	subRoutes []interfaces.Router
// }

// func NewPodsRouterFactory(handlers []func(*gin.RouterGroup) interfaces.Router) func(*gin.RouterGroup) interfaces.Router {
// 	return func(rg *gin.RouterGroup) interfaces.Router {
// 		route := rg.Group("/pods")

// 		var subRoutes []interfaces.Router
// 		for _, handler := range handlers {
// 			subRoutes = append(subRoutes, handler(route))
// 		}

// 		return &podsRouter{
// 			auth:      rg.Group("/pods", m.GetMiddleware().Auth()),
// 			pods:      c.NewPodsController(),
// 			subRoutes: subRoutes,
// 		}
// 	}
// }

// func (c *podsRouter) Init() {
// 	c.auth.POST("/:name", c.pods.Exec)

// 	c.auth.GET("", c.pods.ReadAll)
// 	c.auth.GET("/:namespace", c.pods.ReadAll)
// 	c.auth.GET("/:namespace/:name", c.pods.ReadOne)

// 	for _, route := range c.subRoutes {
// 		route.Init()
// 	}
// }

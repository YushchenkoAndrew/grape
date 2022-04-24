package k3s

import (
	c "api/controllers/k3s"
	"api/interfaces"
	i "api/interfaces/controller"
	m "api/middleware"
	s "api/service/k3s"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

type serviceRouter struct {
	auth    *gin.RouterGroup
	service i.Default
}

func NewServiceRouterFactory(k3s *kubernetes.Clientset) func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &serviceRouter{
			auth:    rg.Group("/service", m.GetMiddleware().Auth()),
			service: c.NewServiceController(s.NewServiceService(k3s)),
		}
	}
}

func (c *serviceRouter) Init() {
	c.auth.POST("/:namespace", c.service.CreateOne)
	c.auth.POST("/list/:namespace", c.service.CreateAll)

	c.auth.GET("/:namespace/:label", c.service.ReadOne)
	c.auth.GET("/:namespace", c.service.ReadAll)

	c.auth.PUT("/:namespace", c.service.UpdateOne)
	c.auth.PUT("/list/:namespace", c.service.UpdateAll)

	c.auth.DELETE("/:namespace/:name", c.service.DeleteAll)
	c.auth.DELETE("/:namespace", c.service.DeleteOne)
}

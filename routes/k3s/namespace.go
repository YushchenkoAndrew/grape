package k3s

import (
	c "api/controllers/k3s"
	"api/interfaces"
	m "api/middleware"
	s "api/service/k3s"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

type namespaceRouter struct {
	auth      *gin.RouterGroup
	namespace interfaces.Default
}

func NewNamespaceRouterFactory(k3s *kubernetes.Clientset) func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &namespaceRouter{
			auth:      rg.Group("/namespace", m.GetMiddleware().Auth()),
			namespace: c.NewNamespaceController(s.NewNamespaceService(k3s)),
		}
	}
}

func (c *namespaceRouter) Init() {
	c.auth.POST("", c.namespace.CreateOne)
	c.auth.POST("/list", c.namespace.CreateAll)

	c.auth.GET("/:label", c.namespace.ReadOne)
	c.auth.GET("", c.namespace.ReadAll)

	c.auth.PUT("", c.namespace.UpdateOne)
	c.auth.PUT("/list", c.namespace.UpdateAll)

	c.auth.DELETE("/:name", c.namespace.DeleteOne)
	c.auth.DELETE("", c.namespace.DeleteAll)
}

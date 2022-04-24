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

type ingressRouter struct {
	auth    *gin.RouterGroup
	ingress i.Default
}

func NewIngressRouterFactory(k3s *kubernetes.Clientset) func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &ingressRouter{
			auth:    rg.Group("/ingress", m.GetMiddleware().Auth()),
			ingress: c.NewIngressController(s.NewIngressService(k3s)),
		}
	}
}

func (c *ingressRouter) Init() {
	c.auth.POST("/:namespace", c.ingress.CreateOne)
	c.auth.POST("/list/:namespace", c.ingress.CreateAll)

	c.auth.GET("/:namespace/:label", c.ingress.ReadOne)
	c.auth.GET("/:namespace", c.ingress.ReadAll)

	c.auth.PUT("/:namespace", c.ingress.UpdateOne)
	c.auth.PUT("/list/:namespace", c.ingress.UpdateAll)

	c.auth.DELETE("/:namespace/:name", c.ingress.DeleteAll)
	c.auth.DELETE("/:namespace", c.ingress.DeleteOne)
}

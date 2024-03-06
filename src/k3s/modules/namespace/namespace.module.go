package namespace

// import (
// 	c "grape/controllers/k3s"
// 	"grape/interfaces"
// 	i "grape/interfaces/controller"
// 	m "grape/middleware"
// 	s "grape/service/k3s"

// 	"github.com/gin-gonic/gin"
// 	"k8s.io/client-go/kubernetes"
// )

// type namespaceRouter struct {
// 	auth      *gin.RouterGroup
// 	namespace i.Default
// }

// func NewNamespaceRouterFactory(k3s *kubernetes.Clientset) func(*gin.RouterGroup) interfaces.Router {
// 	return func(rg *gin.RouterGroup) interfaces.Router {
// 		return &namespaceRouter{
// 			auth:      rg.Group("/namespace", m.GetMiddleware().Auth()),
// 			namespace: c.NewNamespaceController(s.NewNamespaceService(k3s)),
// 		}
// 	}
// }

// func (c *namespaceRouter) Init() {
// 	c.auth.POST("", c.namespace.CreateOne)
// 	c.auth.POST("/list", c.namespace.CreateAll)

// 	c.auth.GET("/:label", c.namespace.ReadOne)
// 	c.auth.GET("", c.namespace.ReadAll)

// 	c.auth.PUT("", c.namespace.UpdateOne)
// 	c.auth.PUT("/list", c.namespace.UpdateAll)

// 	c.auth.DELETE("/:name", c.namespace.DeleteOne)
// 	c.auth.DELETE("", c.namespace.DeleteAll)
// }

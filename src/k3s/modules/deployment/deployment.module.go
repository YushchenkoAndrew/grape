package deployment

// import (
// 	c "grape/controllers/k3s"
// 	"grape/interfaces"
// 	i "grape/interfaces/controller"
// 	m "grape/middleware"
// 	s "grape/service/k3s"

// 	"github.com/gin-gonic/gin"
// 	"k8s.io/client-go/kubernetes"
// )

// type deploymentRouter struct {
// 	auth       *gin.RouterGroup
// 	deployment i.Default
// }

// func NewDeploymentRouterFactory(k3s *kubernetes.Clientset) func(*gin.RouterGroup) interfaces.Router {
// 	return func(rg *gin.RouterGroup) interfaces.Router {
// 		return &deploymentRouter{
// 			auth:       rg.Group("/deployment", m.GetMiddleware().Auth()),
// 			deployment: c.NewDeploymentController(s.NewDeploymentService(k3s)),
// 		}
// 	}
// }

// func (c *deploymentRouter) Init() {
// 	c.auth.POST("/:namespace", c.deployment.CreateOne)
// 	c.auth.POST("/list/:namespace", c.deployment.CreateAll)

// 	c.auth.GET("/:namespace/:label", c.deployment.ReadOne)
// 	c.auth.GET("/:namespace", c.deployment.ReadAll)

// 	c.auth.PUT("/:namespace", c.deployment.UpdateOne)
// 	c.auth.PUT("/list/:namespace", c.deployment.UpdateAll)

// 	c.auth.DELETE("/:namespace/:name", c.deployment.DeleteAll)
// 	c.auth.DELETE("/:namespace", c.deployment.DeleteOne)
// }

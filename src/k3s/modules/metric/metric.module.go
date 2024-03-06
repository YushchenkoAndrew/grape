package metric

// import (
// 	c "grape/controllers/k3s/pods"
// 	"grape/interfaces"
// 	i "grape/interfaces/controller"
// 	m "grape/middleware"
// 	"grape/models"
// 	s "grape/service/k3s/pods"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-redis/redis/v8"
// 	"gorm.io/gorm"
// 	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
// )

// type metricsRouter struct {
// 	auth      *gin.RouterGroup
// 	authToken *gin.RouterGroup
// 	metrics   i.Default
// }

// func NewMetricsRouterFactory(db *gorm.DB, client *redis.Client, metrics *metrics.Clientset) func(*gin.RouterGroup) interfaces.Router {
// 	return func(rg *gin.RouterGroup) interfaces.Router {
// 		return &metricsRouter{
// 			auth:      rg.Group("/metrics", m.GetMiddleware().Auth()),
// 			authToken: rg.Group("/metrics", m.GetMiddleware().AuthToken("SUBSCRIPTION", &[]models.Subscription{})),
// 			metrics:   c.NewMetricsController(s.NewFullMetricsService(db, client, metrics)),
// 		}
// 	}
// }

// func (c *metricsRouter) Init() {
// 	c.authToken.POST("/:namespace/:label", c.metrics.CreateOne)
// 	c.authToken.POST("/:namespace", c.metrics.CreateAll)

// 	c.auth.GET("/:id", c.metrics.ReadOne)
// 	c.auth.GET("", c.metrics.ReadAll)

// 	c.auth.PUT("/:id", c.metrics.UpdateOne)
// 	c.auth.PUT("", c.metrics.UpdateAll)

// 	c.auth.DELETE("/:id", c.metrics.DeleteOne)
// 	c.auth.DELETE("", c.metrics.DeleteAll)
// }

package pods

import (
	c "api/controllers/k3s/pods"
	"api/interfaces"
	m "api/middleware"
	"api/models"
	s "api/service/k3s/pods"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type metricsRouter struct {
	auth      *gin.RouterGroup
	authToken *gin.RouterGroup
	metrics   interfaces.Default
}

func NewMetricsRouterFactory(db *gorm.DB, client *redis.Client, metrics *metrics.Clientset) func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &metricsRouter{
			auth:      rg.Group("/metrics", m.GetMiddleware().Auth()),
			authToken: rg.Group("/metrics", m.GetMiddleware().AuthToken("SUBSCRIPTION", &[]models.Subscription{})),
			metrics:   c.NewMetricsController(s.NewFullMetricsService(db, client, metrics)),
		}
	}
}

func (c *metricsRouter) Init() {
	c.auth.GET("", c.metrics.ReadAll)
	c.auth.GET("/:id", c.metrics.ReadOne)
	// c.auth.GET("/:id/:namespace/:name", c.metrics.ReadOne)

	c.authToken.POST("/:id/:namespace", c.metrics.CreateAll)
	c.authToken.POST("/:id/:namespace/:name", c.metrics.CreateOne)
}

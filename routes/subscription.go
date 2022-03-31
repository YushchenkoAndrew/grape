package routes

import (
	c "api/controllers"
	"api/interfaces"
	m "api/middleware"
	s "api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type subscriptionRouter struct {
	auth         *gin.RouterGroup
	subscription interfaces.Default
}

func NewSubscribeRouter(rg *gin.RouterGroup, db *gorm.DB, client *redis.Client) interfaces.Router {
	return &subscriptionRouter{
		auth:         rg.Group("/subscription", m.GetMiddleware().Auth()),
		subscription: c.NewSubscriptionController(s.NewFullSubscriptionService(db, client)),
	}
}

func (c *subscriptionRouter) Init() {
	c.auth.POST("/:id", c.subscription.CreateOne)
	c.auth.POST("/list/:id", c.subscription.CreateAll)

	c.auth.GET("/:id", c.subscription.ReadOne)
	c.auth.GET("", c.subscription.ReadAll)

	c.auth.PUT("/:id", c.subscription.UpdateOne)
	c.auth.PUT("", c.subscription.UpdateAll)

	c.auth.DELETE("/:id", c.subscription.DeleteOne)
	c.auth.DELETE("", c.subscription.DeleteAll)
}

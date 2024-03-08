package statistic

import (
	"grape/src/common/middleware"
	m "grape/src/common/module"

	"github.com/gin-gonic/gin"
)

type infoModule struct {
	*m.Module[StatisticT]
}

func NewInfoModule(rg *gin.RouterGroup, handlers []func(*gin.RouterGroup) m.ModuleT) m.ModuleT {
	// var subRoutes []interfaces.Router
	// for _, handler := range handlers {
	// 	subRoutes = append(subRoutes, handler(route))
	// }

	// TODO:
	return &infoModule{
		Module: &m.Module[StatisticT]{
			Route: rg.Group("/info"),
			Auth:  rg.Group("/info", middleware.GetMiddleware().Jwt()),
			// info:      c.NewInfoController(),
			// SubRoutes: subRoutes,
		},
	}
}

// FIXME: !!!
func (c *infoModule) Init() {
	// c.auth.POST("", c.info.Create)
	// c.auth.POST("/list", c.info.CreateAll)
	// c.auth.POST("/:date", c.info.CreateOne)

	// c.route.GET("", c.info.ReadAll)
	// c.route.GET("/:id", c.info.ReadOne)

	// c.auth.PUT("", c.info.UpdateAll)
	// c.auth.PUT("/:id", c.info.UpdateOne)

	// c.auth.DELETE("", c.info.DeleteAll)
	// c.auth.DELETE("/:id", c.info.DeleteOne)

	// for _, route := range c.subRoutes {
	// 	route.Init()
	// }
}

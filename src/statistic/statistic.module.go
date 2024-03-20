package statistic

import (
	"grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type statisticModule struct {
	*m.Module[*StatisticController]
}

func NewStatisticModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &statisticModule{
		Module: &m.Module[*StatisticController]{
			Route:      rg.Group("/statistics"),
			Auth:       rg.Group("/statistics", middleware.GetMiddleware(nil).Jwt()),
			Controller: NewStatisticController(NewStatisticService(s)),
			Modules:    modules,
		},
	}
}

func (c *statisticModule) Init() {
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

package context

import (
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type contextModule struct {
	*m.Module[*ContextController]
}

func NewContextModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &contextModule{
		Module: &m.Module[*ContextController]{
			Route:      rg.Group("/contexts"),
			Auth:       rg.Group("/admin/contexts", h.GetMiddleware(nil).Jwt()),
			Controller: NewContextController(NewContextService(s)),
			Modules:    modules,
		},
	}
}

func (c *contextModule) Init() {
	// c.Route.GET("/:id", c.Controller.FindOne)

	// c.Auth.GET("/:id", c.Controller.AdminFindOne)

	// c.Auth.POST("", c.Controller.Create)

	// c.Auth.PUT("/:id", c.Controller.Update)
	// c.Auth.PUT("/:id/order", c.Controller.PutOrder)

	// c.Auth.DELETE("/:id", c.Controller.Delete)

	c.Module.Init()
}

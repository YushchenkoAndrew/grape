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

	c.Auth.GET("/:id", c.Controller.AdminFindOne)

	c.Auth.POST("", c.Controller.Create)
	c.Auth.POST("/:id/fields", c.Controller.CreateField)

	c.Auth.PUT("/:id", c.Controller.Update)
	c.Auth.PUT("/:id/order", c.Controller.UpdateOrder)
	c.Auth.PUT("/:id/fields/:field_id", c.Controller.UpdateField)
	c.Auth.PUT("/:id/fields/:field_id/order", c.Controller.UpdateFieldOrder)

	c.Auth.DELETE("/:id", c.Controller.Delete)
	c.Auth.DELETE("/:id/fields/:field_id", c.Controller.DeleteField)

	c.Module.Init()
}

package link

import (
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type linkModule struct {
	*m.Module[*LinkController]
}

func NewLinkModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &linkModule{
		Module: &m.Module[*LinkController]{
			Route:      rg.Group("/links"),
			Auth:       rg.Group("/admin/links", h.GetMiddleware(nil).Jwt()),
			Controller: NewLinkController(NewLinkService(s)),
			Modules:    modules,
		},
	}
}

func (c *linkModule) Init() {
	c.Auth.GET("/:id", c.Controller.AdminFindOne)

	c.Auth.POST("", c.Controller.Create)

	c.Auth.PUT("/:id", c.Controller.Update)
	c.Auth.PUT("/:id/order", c.Controller.UpdateOrder)

	c.Auth.DELETE("/:id", c.Controller.Delete)

	c.Module.Init()
}

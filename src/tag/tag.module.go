package tag

import (
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type tagModule struct {
	*m.Module[*TagController]
}

func NewTagModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &tagModule{
		Module: &m.Module[*TagController]{
			Route:      rg.Group("/tags"),
			Auth:       rg.Group("/admin/tags", h.GetMiddleware(nil).Jwt()),
			Controller: NewTagController(NewTagService(s)),
			Modules:    modules,
		},
	}
}

func (c *tagModule) Init() {
	c.Auth.GET("/:id", c.Controller.AdminFindOne)

	c.Auth.POST("", c.Controller.Create)

	c.Auth.PUT("/:id", c.Controller.Update)

	c.Auth.DELETE("/:id", c.Controller.Delete)

	c.Module.Init()
}

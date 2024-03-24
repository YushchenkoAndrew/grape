package palette

import (
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type module struct {
	*m.Module[*PaletteController]
}

func NewPaletteModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &module{
		Module: &m.Module[*PaletteController]{
			Route:      rg.Group("/palettes"),
			Controller: NewPaletteController(NewPaletteService(s)),
			Modules:    modules,
		},
	}
}

func (c *module) Init() {
	c.Route.GET("", c.Controller.FindAll)
	c.Route.GET("/:id", c.Controller.FindOne)

	c.Route.POST("", c.Controller.Create)
	c.Route.PUT("/:id", c.Controller.Update)
	c.Route.DELETE("/:id", c.Controller.Delete)

	c.Module.Init()
}

package filter

import (
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type filterModule struct {
	*m.Module[FilterT]
}

func NewFilterModule(rg *gin.RouterGroup, modules *[]m.ModuleT, client *service.CommonService) m.ModuleT {
	return &filterModule{
		Module: &m.Module[FilterT]{
			Route:      rg,
			Controller: NewFilterController(NewFilterService(client)),
			Modules:    modules,
		},
	}
}

func (c *filterModule) Init() {
	c.Route.GET("/trace/:ip", c.Controller.TraceIp)
	c.Module.Init()
}

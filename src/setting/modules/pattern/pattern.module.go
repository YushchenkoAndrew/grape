package pattern

import (
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type module struct {
	*m.Module[*PatternController]
}

func NewPatternModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &module{
		Module: &m.Module[*PatternController]{
			Route:      rg.Group("/pattern"),
			Controller: NewPatternController(NewPatternService(s)),
			Modules:    modules,
		},
	}
}

func (c *module) Init() {
	c.Module.Init()
}

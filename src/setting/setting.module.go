package setting

import (
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/setting/modules/palette"

	"github.com/gin-gonic/gin"
)

type module struct {
	*m.Module[*SettingController]
}

func NewSettingModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	group := rg.Group("/admin/settings", h.GetMiddleware(nil).Jwt())

	return &module{
		Module: &m.Module[*SettingController]{
			Route:      rg.Group("/settings"),
			Auth:       group,
			Controller: NewSettingController(NewSettingService(s)),
			Modules: append(modules,
				palette.NewPaletteModule(group, []m.ModuleT{}, s),
			),
		},
	}
}

func (c *module) Init() {
	c.Module.Init()
}

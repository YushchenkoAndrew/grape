package project

import (
	c "grape/src/common/controller"
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type projectModule struct {
	*m.Module[c.CommonController]
}

func NewProjectModule(rg *gin.RouterGroup, modules *[]m.ModuleT, s *service.CommonService) m.ModuleT {
	return &projectModule{
		Module: &m.Module[c.CommonController]{
			Route:      rg.Group("/project"),
			Auth:       rg.Group("/admin/project", h.GetMiddleware(nil).Jwt()),
			Controller: NewProjectController(NewProjectService(s)),
		},
	}
}

func (c *projectModule) Init() {
	c.Route.GET("", c.Controller.FindAll)
	c.Route.GET("/:id", c.Controller.FindOne)

	c.Auth.GET("", c.Controller.FindAll)
	c.Auth.GET("/:id", c.Controller.FindOne)

	c.Auth.POST("", c.Controller.Create)
	c.Auth.PUT("/:id", c.Controller.Update)
	c.Auth.DELETE("/:id", c.Controller.Delete)

	c.Module.Init()
}

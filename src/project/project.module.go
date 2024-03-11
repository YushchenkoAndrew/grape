package project

import (
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type projectModule struct {
	*m.Module[*ProjectController]
}

func NewProjectModule(rg *gin.RouterGroup, modules *[]m.ModuleT, s *service.CommonService) m.ModuleT {
	return &projectModule{
		Module: &m.Module[*ProjectController]{
			Route:      rg.Group("/projects"),
			Auth:       rg.Group("/admin/projects", h.GetMiddleware(nil).Jwt()),
			Controller: NewProjectController(NewProjectService(s)),
		},
	}
}

func (c *projectModule) Init() {
	c.Route.GET("", c.Controller.FindAll)
	c.Route.GET("/:id", c.Controller.FindOne)

	c.Auth.GET("", c.Controller.AdminFindAll)
	c.Auth.GET("/:id", c.Controller.AdminFindOne)

	c.Auth.POST("", c.Controller.Create)
	c.Auth.PUT("/:id", c.Controller.Update)
	c.Auth.DELETE("/:id", c.Controller.Delete)

	c.Module.Init()
}

package project

import (
	"grape/src/common/client"
	c "grape/src/common/controller"
	"grape/src/common/middleware"
	m "grape/src/common/module"

	"github.com/gin-gonic/gin"
)

type projectModule struct {
	*m.Module[c.DefaultController]
}

func NewProjectModule(rg *gin.RouterGroup, modules *[]m.ModuleT, client *client.Clients) m.ModuleT {
	return &projectModule{
		Module: &m.Module[c.DefaultController]{
			Route:      rg.Group("/project"),
			Auth:       rg.Group("/project", middleware.GetMiddleware().Auth()),
			Controller: NewProjectController(NewProjectService(client)),
		},
	}
}

func (c *projectModule) Init() {
	c.Route.GET("", c.Controller.FindAll)
	c.Route.GET("/:name", c.Controller.FindOne)

	c.Auth.POST("", c.Controller.Create)
	c.Auth.PUT("/:name", c.Controller.Update)
	c.Auth.DELETE("/:name", c.Controller.Delete)

	c.Module.Init()
}

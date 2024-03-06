package src

import (
	"grape/src/common/client"
	m "grape/src/common/module"

	"github.com/gin-gonic/gin"
)

type indexModule struct {
	*m.Module[IndexT]
}

func NewIndexModule(route *gin.RouterGroup, modules *[]m.ModuleT, client *client.Clients) m.ModuleT {
	return &indexModule{
		Module: &m.Module[IndexT]{
			Route:      route,
			Controller: NewIndexController(NewIndexService(client)),
			Modules:    modules,
		},
	}
}

func (c *indexModule) Init() {
	c.Route.GET("/ping", c.Controller.Ping)

	c.Module.Init()
}

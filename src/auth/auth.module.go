package auth

import (
	"grape/src/common/client"
	m "grape/src/common/module"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	*m.Module[AuthT]
}

func NewAuthRouter(route *gin.RouterGroup, modules *[]m.ModuleT, client *client.Clients) m.ModuleT {
	return &authRouter{
		Module: &m.Module[AuthT]{
			Route:      route,
			Controller: NewAuthController(NewAuthService(client)),
			Modules:    modules,
		},
	}
}

func (c *authRouter) Init() {
	c.Route.POST("/login", c.Controller.Login)
	c.Route.POST("/refresh", c.Controller.Refresh)

	c.Module.Init()
}

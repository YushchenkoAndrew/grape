package auth

import (
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	*m.Module[*AuthController]
}

func NewAuthModule(route *gin.RouterGroup, modules *[]m.ModuleT, s *service.CommonService) m.ModuleT {
	return &authRouter{
		Module: &m.Module[*AuthController]{
			Route:      route,
			Controller: NewAuthController(NewAuthService(s)),
			Modules:    modules,
		},
	}
}

func (c *authRouter) Init() {
	c.Route.POST("/login", c.Controller.Login)
	c.Route.POST("/refresh", c.Controller.Refresh)

	c.Module.Init()
}

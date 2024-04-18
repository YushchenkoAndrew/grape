package auth

import (
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type authRouter struct {
	*m.Module[*AuthController]
}

func NewAuthModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &authRouter{
		Module: &m.Module[*AuthController]{
			Route:      rg.Group("/auth"),
			Auth:       rg.Group("/auth", h.GetMiddleware(nil).Jwt()),
			Controller: NewAuthController(NewAuthService(s)),
			Modules:    modules,
		},
	}
}

func (c *authRouter) Init() {
	c.Route.POST("/login", c.Controller.Login)
	c.Route.POST("/refresh", c.Controller.Refresh)
	c.Auth.POST("/logout", c.Controller.Logout)
	c.Auth.GET("/ping", c.Controller.Ping)

	c.Module.Init()
}

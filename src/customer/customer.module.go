package customer

import (
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type customerModule struct {
	*m.Module[*CustomerController]
}

func NewCustomerModule(rg *gin.RouterGroup, modules *[]m.ModuleT, client *service.CommonService) m.ModuleT {
	return &customerModule{
		Module: &m.Module[*CustomerController]{
			Route:      rg,
			Controller: NewCustomerController(NewCustomerService(client)),
			Modules:    modules,
		},
	}
}

func (c *customerModule) Init() {
	c.Route.GET("/trace/:ip", c.Controller.TraceIp)
	c.Module.Init()
}

package swagger

import (
	_ "grape/docs"
	m "grape/src/common/module"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type swaggerModule struct {
	*m.Module[interface{}]
}

func NewSwaggerRouter(route *gin.RouterGroup) m.ModuleT {
	return &swaggerModule{
		Module: &m.Module[interface{}]{Route: route.Group("/swagger")},
	}
}

func (c *swaggerModule) Init() {
	c.Route.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

package attachment

import (
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type attachmentModule struct {
	*m.Module[*AttachmentController]
}

func NewAttachmentModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &attachmentModule{
		Module: &m.Module[*AttachmentController]{
			Route:      rg.Group("/attachments"),
			Auth:       rg.Group("/admin/attachments", h.GetMiddleware(nil).Jwt()),
			Controller: NewAttachmentController(NewAttachmentService(s)),
			Modules:    modules,
		},
	}
}

func (c *attachmentModule) Init() {
	c.Route.GET("/:id", c.Controller.FindOne)

	c.Auth.GET("/:id", c.Controller.AdminFindOne)

	c.Auth.POST("", c.Controller.Create)

	c.Auth.PUT("/:id", c.Controller.Update)
	c.Auth.PUT("/:id/order", c.Controller.PutOrder)

	c.Auth.DELETE("/:id", c.Controller.Delete)

	c.Module.Init()
}

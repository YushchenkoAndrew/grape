package stage

import (
	h "grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"

	"github.com/gin-gonic/gin"
)

type projectModule struct {
	*m.Module[*StageController]
}

func NewStageModule(rg *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
	return &projectModule{
		Module: &m.Module[*StageController]{
			Route:      rg.Group("/stages"),
			Auth:       rg.Group("/admin/stages", h.GetMiddleware(nil).Jwt()),
			Controller: NewStageController(NewStageService(s)),
			Modules:    modules,
		},
	}
}

func (c *projectModule) Init() {
	c.Route.GET("", c.Controller.FindAll)

	c.Auth.GET("", c.Controller.AdminFindAll)

	c.Auth.POST("", c.Controller.Create)
	c.Auth.POST("/:id/tasks", c.Controller.CreateTask)

	c.Auth.PUT("/:id", c.Controller.Update)
	c.Auth.PUT("/:id/order", c.Controller.UpdateOrder)
	c.Auth.PUT("/:id/tasks/:task_id", c.Controller.UpdateTask)
	c.Auth.PUT("/:id/tasks/:task_id/order", c.Controller.UpdateTaskOrder)

	c.Auth.DELETE("/:id", c.Controller.Delete)
	c.Auth.DELETE("/:id/tasks/:task_id", c.Controller.DeleteTask)

	c.Module.Init()
}

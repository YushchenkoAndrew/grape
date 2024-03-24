package module

import (
	"github.com/gin-gonic/gin"
)

type ModuleT interface {
	Init()
}

type Module[T any] struct {
	Route      *gin.RouterGroup
	Auth       *gin.RouterGroup
	Controller T
	Modules    []ModuleT
}

func (c *Module[T]) Init() {
	if c.Modules == nil {
		return
	}

	for _, route := range c.Modules {
		route.Init()
	}
}

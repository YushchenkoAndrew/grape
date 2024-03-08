package test

import (
	"grape/src/common/client"
	"grape/src/common/config"
	"grape/src/common/middleware"
	"grape/src/common/service"

	m "grape/src/common/module"
	_ "grape/src/common/validator"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(module func(route *gin.RouterGroup, modules *[]m.ModuleT, s *service.CommonService) m.ModuleT) *gin.Engine {
	cfg := config.NewConfig("configs/", "config", "yaml")

	service := &service.CommonService{
		DB:     client.ConnPsql(cfg),
		Redis:  client.ConnRedis(cfg),
		Config: cfg,
	}

	r := gin.Default()
	rg := r.Group(cfg.Server.Prefix, middleware.NewMiddleware(service).Limit())
	module(rg, &[]m.ModuleT{}, service).Init()
	return r
}

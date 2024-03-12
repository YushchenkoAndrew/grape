package main

import (
	"fmt"
	"grape/src"
	"grape/src/auth"
	"grape/src/common/client"
	"grape/src/common/config"
	"grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/filter"
	"grape/src/project"
	"grape/src/swagger"

	_ "grape/src/common/validator"

	"github.com/gin-gonic/gin"
)

// @title Gin API
// @version 1.0
// @description Remake of my previous attampted on creating API with Node.js

// @contact.name API Author
// @contact.url https://mortis-grimreaper.ddns.net/projects
// @contact.email AndrewYushchenko@gmail.com

// @license.name MIT
// @license.url https://github.com/YushchenkoAndrew/API_Server/blob/master/LICENSE

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host mortis-grimreaper.ddns.net
// @BasePath /api
func main() {
	// config.NewConfigs([]func() config.ConfigT{
	// 	config.NewEnvConfig("./", ""),
	// 	// config.NewOperationConfig("./", "operations"),
	// }).Init()

	// cfg
	cfg := config.NewConfig("configs/", "config", "yaml")

	service := &service.CommonService{
		DB:     client.ConnPsql(cfg),
		Redis:  client.ConnRedis(cfg),
		Config: cfg,
	}

	r := gin.Default()
	rg := r.Group(cfg.Server.Prefix, middleware.GetMiddleware(service).Default())
	module := src.NewIndexModule(rg, &[]m.ModuleT{
		swagger.NewSwaggerRouter(rg),
		auth.NewAuthModule(rg, &[]m.ModuleT{}, service),
		filter.NewFilterModule(rg, &[]m.ModuleT{}, service),
		project.NewProjectModule(rg, &[]m.ModuleT{}, service),
		// routes.NewFileRouter(rg, db, redis),
		// routes.NewLinkRouter(rg, db, redis),
		// routes.NewBotRouter(rg, db, redis),
		// routes.NewPatternRouter(rg, db, redis),

		// // routes.NewWorldRouter(rg),
		// // routes.NewInfoRouter(rg, []func(*gin.RouterGroup) interfaces.Router{
		// // 	info.NewSumRouterFactory(),
		// // 	info.NewRangeRouterFactory(),
		// // }),

		// routes.NewK3sRouter(rg, []func(*gin.RouterGroup) interfaces.Router{
		// 	k.NewDeploymentRouterFactory(k3s),
		// 	k.NewIngressRouterFactory(k3s),
		// 	k.NewPodsRouterFactory([]func(*gin.RouterGroup) interfaces.Router{
		// 		pods.NewMetricsRouterFactory(db, redis, metrics),
		// 	}),
		// 	k.NewNamespaceRouterFactory(k3s),
		// 	k.NewServiceRouterFactory(k3s),
		// }),

		// routes.NewSubscribeRouter(rg, db, redis),
	}, service)

	module.Init()
	r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}

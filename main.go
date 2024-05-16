package main

import (
	"grape/src"
	"grape/src/attachment"
	"grape/src/auth"
	"grape/src/common/client"
	"grape/src/common/config"
	"grape/src/common/dto/response"
	"grape/src/common/middleware"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/context"
	"grape/src/customer"
	"grape/src/link"
	"grape/src/project"
	"grape/src/setting"
	"grape/src/stage"
	"grape/src/statistic"
	"grape/src/swagger"
	"grape/src/tag"
	"net/http"

	_ "grape/src/common/validator"

	"github.com/gin-gonic/gin"
)

// @title Gin API
// @version 3.0
// @description Gin API for my various projects

// @contact.name API Author
// @contact.email AndrewYushchenko@gmail.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @BasePath /grape
func main() {
	cfg := config.NewConfig("configs/", "config", "yaml")

	service := &service.CommonService{
		DB:     client.ConnPsql(cfg),
		Redis:  client.ConnRedis(cfg),
		Config: cfg,
	}

	r := gin.Default()
	rg := r.Group(cfg.Server.Prefix, middleware.GetMiddleware(service).Default())
	module := src.NewIndexModule(rg, []m.ModuleT{
		swagger.NewSwaggerRouter(rg),
		auth.NewAuthModule(rg, []m.ModuleT{}, service),
		customer.NewCustomerModule(rg, []m.ModuleT{}, service),
		project.NewProjectModule(rg, []m.ModuleT{}, service),
		attachment.NewAttachmentModule(rg, []m.ModuleT{}, service),
		link.NewLinkModule(rg, []m.ModuleT{}, service),
		statistic.NewStatisticModule(rg, []m.ModuleT{}, service),
		setting.NewSettingModule(rg, []m.ModuleT{}, service),
		stage.NewStageModule(rg, []m.ModuleT{}, service),
		context.NewContextModule(rg, []m.ModuleT{}, service),
		tag.NewTagModule(rg, []m.ModuleT{}, service),
		// routes.NewBotRouter(rg, db, redis),

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

	r.NoRoute(func(ctx *gin.Context) { response.ThrowErr(ctx, http.StatusNotFound, "Resource not found") })
	r.Run(cfg.Server.Address)
}

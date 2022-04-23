package main

import (
	"api/client"
	"api/config"
	"api/interfaces"
	m "api/middleware"
	"api/models"
	"api/routes"
	k "api/routes/k3s"
	"api/routes/k3s/pods"

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
	config.NewConfig([]func() interfaces.Config{
		config.NewEnvConfig("./", ""),
		config.NewOperationConfig("./", "operations"),
	}).Init()

	redis := client.ConnRedis()
	k3s, metrics := client.ConnK3s()
	db := client.ConnDB([]interfaces.Table{
		models.NewInfo(),
		models.NewWorld(),

		models.NewGeoIpBlocks(),
		models.NewGeoIpLocations(),
		models.NewPattern(),

		models.NewFile(),
		models.NewLink(),
		models.NewMetrics(),
		models.NewSubscription(),
		models.NewProject(),
	})

	r := gin.Default()
	rg := r.Group(config.ENV.BasePath, m.NewMiddleware(db, redis).Limit())
	router := routes.NewIndexRouter(rg, &[]interfaces.Router{
		routes.NewSwaggerRouter(rg),
		routes.NewProjectRouter(rg, db, redis),
		routes.NewFileRouter(rg, db, redis),
		routes.NewLinkRouter(rg, db, redis),
		routes.NewBotRouter(rg, db, redis),
		routes.NewPatternRouter(rg, db, redis),

		// routes.NewWorldRouter(rg),
		// routes.NewInfoRouter(rg, []func(*gin.RouterGroup) interfaces.Router{
		// 	info.NewSumRouterFactory(),
		// 	info.NewRangeRouterFactory(),
		// }),

		routes.NewK3sRouter(rg, []func(*gin.RouterGroup) interfaces.Router{
			k.NewDeploymentRouterFactory(k3s),
			k.NewIngressRouterFactory(k3s),
			k.NewPodsRouterFactory([]func(*gin.RouterGroup) interfaces.Router{
				pods.NewMetricsRouterFactory(db, redis, metrics),
			}),
			k.NewNamespaceRouterFactory(k3s),
			k.NewServiceRouterFactory(k3s),
		}),

		routes.NewSubscribeRouter(rg, db, redis),
	}, db, redis)

	router.Init()
	r.Run(config.ENV.Host + ":" + config.ENV.Port)
}

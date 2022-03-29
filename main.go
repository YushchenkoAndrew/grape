package main

import (
	"api/config"
	"api/db"
	"api/interfaces"
	m "api/middleware"
	"api/models"
	"api/routes"

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

// @host mortis-grimreaper.ddns.net:31337
// @BasePath /api
func main() {
	config.NewConfig([]func() interfaces.Config{
		config.NewEnvConfig("./"),
		config.NewK3sConfig("./k3s.yaml"),
		config.NewOperationConfig("./", "operations"),
	}).Init()

	db, client := db.Init([]interfaces.Table{
		models.NewInfo(),
		models.NewWorld(),

		models.NewGeoIpBlocks(),
		models.NewGeoIpLocations(),

		models.NewFile(),
		models.NewLink(),
		models.NewMetrics(),
		models.NewSubscription(),
		models.NewProject(),
	})

	r := gin.Default()
	rg := r.Group(config.ENV.BasePath, m.NewMiddleware(client).Limit())
	router := routes.NewIndexRouter(rg, &[]interfaces.Router{
		routes.NewSwaggerRouter(rg),
		// routes.NewWorldRouter(rg),
		routes.NewProjectRouter(rg, db, client),
		routes.NewFileRouter(rg, db, client),
		routes.NewLinkRouter(rg, db, client),
		// routes.NewBotRouter(rg),

		// routes.NewInfoRouter(rg, []func(*gin.RouterGroup) interfaces.Router{
		// 	info.NewSumRouterFactory(),
		// 	info.NewRangeRouterFactory(),
		// }),

		// routes.NewK3sRouter(rg, []func(*gin.RouterGroup) interfaces.Router{
		// 	k3s.NewDeploymentRouterFactory(),
		// 	k3s.NewIngressRouterFactory(),
		// 	k3s.NewPodsRouterFactory([]func(*gin.RouterGroup) interfaces.Router{
		// 		pods.NewMetricsRouterFactory(),
		// 	}),
		// 	k3s.NewNamespaceRouterFactory(),
		// 	k3s.NewServiceRouterFactory(),
		// }),

		// routes.NewSubscribeRouter(rg),
	})

	router.Init()
	r.Run(config.ENV.Host + ":" + config.ENV.Port)
}

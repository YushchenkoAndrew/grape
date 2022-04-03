package services_test

import (
	"api/config"
	"api/db"
	"api/interfaces"
	m "api/models"
	"api/service"
)

func init() {
	config.NewConfig([]func() interfaces.Config{
		config.NewEnvConfig("./", ""),
	}).Init()

	db, client := db.Init([]interfaces.Table{
		m.NewFile(),
		m.NewLink(),
		m.NewSubscription(),
		m.NewProject(),
	})

	file = *service.NewFileService(db, client)
	link = *service.NewLinkService(db, client)
	subscription = *service.NewSubscriptionService(db, client)
	project = *service.NewProjectService(db, client)

	project.Create(&m.Project{ID: 1, Name: "yes", Title: "js", Flag: "js", Desc: "js", Note: "js"})
}

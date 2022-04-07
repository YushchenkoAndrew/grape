package client

import (
	"api/config"
	"api/interfaces"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	m "k8s.io/metrics/pkg/client/clientset/versioned"
)

func Init(tables []interfaces.Table) (*gorm.DB, *redis.Client, *kubernetes.Clientset, *m.Clientset) {
	var db = connDB()
	for _, table := range tables {
		if err := table.Migrate(db, config.ENV.ForceMigrate); err != nil {
			panic(err)
		}
	}

	var k3s, metrics = connK3s()
	return db, connRedis(), k3s, metrics
}

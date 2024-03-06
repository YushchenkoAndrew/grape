package client

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Clients struct {
	DB    *gorm.DB
	Redis *redis.Client
}

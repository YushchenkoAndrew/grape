package middleware

import (
	"grape/src/common/client"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Middleware struct {
	db     *gorm.DB
	client *redis.Client
}

var middleware Middleware

func NewMiddleware(client *client.Clients) *Middleware {
	middleware = Middleware{db: client.DB, client: client.Redis}
	return &middleware
}

func GetMiddleware() *Middleware {
	return &middleware
}

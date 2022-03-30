package middleware

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Middleware struct {
	db     *gorm.DB
	client *redis.Client
}

var middleware Middleware

func NewMiddleware(db *gorm.DB, client *redis.Client) *Middleware {
	middleware = Middleware{db: db, client: client}
	return &middleware
}

func GetMiddleware() *Middleware {
	return &middleware
}

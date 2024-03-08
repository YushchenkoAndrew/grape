package middleware

import (
	"grape/src/common/service"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Middleware struct {
	db    *gorm.DB
	redis *redis.Client
}

var middleware *Middleware

func NewMiddleware(s *service.CommonService) *Middleware {
	middleware = &Middleware{db: s.DB, redis: s.Redis}
	return middleware
}

func GetMiddleware() *Middleware {
	return middleware
}

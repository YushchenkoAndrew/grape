package middleware

import (
	"grape/src/common/config"
	"grape/src/common/service"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Middleware struct {
	db     *gorm.DB
	redis  *redis.Client
	config *config.Config
}

var middleware *Middleware

func GetMiddleware(s *service.CommonService) *Middleware {
	if middleware != nil || s == nil {
		return middleware
	}

	middleware = &Middleware{db: s.DB, redis: s.Redis, config: s.Config}
	return middleware
}

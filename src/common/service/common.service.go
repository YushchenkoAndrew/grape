package service

import (
	"grape/src/common/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CommonService struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
}

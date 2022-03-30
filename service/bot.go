package service

import (
	m "api/models"
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type BotService struct {
	db     *gorm.DB
	client *redis.Client
}

func NewBotService(db *gorm.DB, client *redis.Client) *BotService {
	return &BotService{db, client}
}

func (c *BotService) RunRedis(dto *m.BotRedisDto) ([]string, error) {
	var command []interface{}
	for _, word := range strings.Split(dto.Command, " ") {
		command = append(command, word)
	}

	return c.client.Do(context.Background(), command...).StringSlice()
}

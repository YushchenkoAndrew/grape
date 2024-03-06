package auth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	m "grape/models"
	"grape/src/common/client"
	"grape/src/common/config"
	"grape/src/common/helper"
	v "grape/src/common/validation"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type authService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAuthService(client *client.Clients) *authService {
	return &authService{db: client.DB, redis: client.Redis}
}

func (c *authService) Login(dto *m.LoginDto) (*m.Auth, error) {
	hasher := md5.New()
	pass := strings.Split(dto.Pass, "$")
	hasher.Write([]byte(pass[0] + config.ENV.Pepper + config.ENV.Pass))

	if !v.ValidateStr(dto.User, config.ENV.User) ||
		!v.ValidateStr(hex.EncodeToString(hasher.Sum(nil)), pass[1]) {
		return nil, fmt.Errorf("Invalid login inforamation")
	}

	var t, err = helper.CreateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	ctx := context.Background()
	token := m.NewAuth().Fill(t)
	c.redis.Set(ctx, token.AccessUUID, config.ENV.ID, time.Duration((token.AccessExpire-now)*int64(time.Second)))
	c.redis.Set(ctx, token.RefreshUUID, config.ENV.ID, time.Duration((token.RefreshExpire-now)*int64(time.Second)))

	return token, nil
}

func (c *authService) Refresh(dto *m.TokenDto) (*m.Auth, error) {
	t, err := helper.CheckToken(dto.RefreshToken)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	auth := m.NewAuth().Fill(t)

	// Double check if such UUID exist in cache + it's the same user
	// (btw don't need it, I have only one user)
	if cacheUUID, err := c.redis.Get(ctx, auth.RefreshUUID).Result(); err != nil || cacheUUID != auth.AccessUUID {
		return nil, fmt.Errorf("Invalid token inforamation")
	}

	c.redis.Del(ctx, auth.RefreshUUID)

	t, err = helper.CreateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	token := m.NewAuth().Fill(t)
	c.redis.Set(ctx, token.AccessUUID, config.ENV.ID, time.Duration((token.AccessExpire-now)*int64(time.Second)))
	c.redis.Set(ctx, token.RefreshUUID, config.ENV.ID, time.Duration((token.RefreshExpire-now)*int64(time.Second)))

	return token, nil
}

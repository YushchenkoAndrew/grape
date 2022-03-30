package service

import (
	"api/config"
	"api/helper"
	m "api/models"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type IndexService struct {
	key string

	db     *gorm.DB
	client *redis.Client
}

func NewIndexService(db *gorm.DB, client *redis.Client) *IndexService {
	return &IndexService{key: "INDEX", db: db, client: client}
}

func (c *IndexService) TraceIP(ip string) ([]m.GeoIpLocations, error) {
	var model []m.GeoIpLocations
	err, _ := helper.Getcache(c.db.Where("geoname_id IN (?)", c.db.Select("geoname_id").Where("network >>= ?::inet", ip).Model(&m.GeoIpBlocks{})), c.client, c.key, fmt.Sprintf("BLOCK:%s", ip), &model)
	return model, err
}

func (c *IndexService) Login(dto *m.LoginDto) (*m.Auth, error) {
	hasher := md5.New()
	pass := strings.Split(dto.Pass, "$")
	hasher.Write([]byte(pass[0] + config.ENV.Pepper + config.ENV.Pass))

	if !helper.ValidateStr(dto.User, config.ENV.User) ||
		!helper.ValidateStr(hex.EncodeToString(hasher.Sum(nil)), pass[1]) {
		return nil, fmt.Errorf("Invalid login inforamation")
	}

	var token m.Auth
	if err := helper.CreateToken(&token); err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	ctx := context.Background()
	c.client.Set(ctx, token.AccessUUID, config.ENV.ID, time.Duration((token.AccessExpire-now)*int64(time.Second)))
	c.client.Set(ctx, token.RefreshUUID, config.ENV.ID, time.Duration((token.RefreshExpire-now)*int64(time.Second)))

	return &token, nil
}

func (c *IndexService) Refresh(dto *m.TokenDto) (*m.Auth, error) {
	auth, err := helper.CheckToken(dto)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Double check if such UUID exist in cache + it's the same user
	// (btw don't need it, I have only one user)
	if cacheUUID, err := c.client.Get(ctx, auth.RefreshUUID).Result(); err != nil || cacheUUID != auth.AccessUUID {
		return nil, fmt.Errorf("Invalid token inforamation")
	}

	c.client.Del(ctx, auth.RefreshUUID)

	var token m.Auth
	if err := helper.CreateToken(&token); err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	c.client.Set(ctx, token.AccessUUID, config.ENV.ID, time.Duration((token.AccessExpire-now)*int64(time.Second)))
	c.client.Set(ctx, token.RefreshUUID, config.ENV.ID, time.Duration((token.RefreshExpire-now)*int64(time.Second)))

	return &token, nil
}

package auth

import (
	"context"
	"fmt"
	"grape/src/auth/dto/request"
	"grape/src/auth/dto/response"
	t "grape/src/auth/types"
	"grape/src/common/config"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/common/types"
	u "grape/src/user/dto/response"
	e "grape/src/user/entities"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	redis  *redis.Client
	config *config.Config
}

func NewAuthService(s *service.CommonService) *AuthService {
	return &AuthService{db: s.DB, redis: s.Redis, config: s.Config}
}

func (c *AuthService) Login(dto *request.LoginDto) (*response.LoginResponseDto, error) {
	var users []e.UserEntity
	if c.db.Joins("Organization").Limit(1).Find(&users, "users.name = ? AND users.status = ?", dto.Username, types.Active); len(users) == 0 {
		return nil, fmt.Errorf("user '%s' not found", dto.Username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(dto.Password)); err != nil {
		return nil, fmt.Errorf("user '%s' not found2", dto.Username)
	}

	return c.generate(&users[0])
}

func (c *AuthService) Refresh(dto *request.RefreshDto) (*response.LoginResponseDto, error) {
	var claim t.RefreshClaim
	token, err := jwt.ParseWithClaims(dto.RefreshToken, &claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(c.config.Jwt.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("unauthorized token")
	}

	ctx := context.Background()
	user_id, err := c.redis.Get(ctx, claim.UID).Result()

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	defer func() {
		c.redis.Del(ctx, claim.UID)
	}()

	var users []e.UserEntity
	if c.db.Joins("Organization").Limit(1).Find(&users, "users.uuid = ? AND users.status = ?", user_id, types.Active); len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return c.generate(&users[0])
}

func (c *AuthService) Logout(claim *t.AccessClaim) (interface{}, error) {

	ctx := context.Background()
	c.redis.Del(ctx, claim.UID)
	c.redis.Del(ctx, claim.RID)

	return nil, nil
}

func (c *AuthService) generate(user *e.UserEntity) (*response.LoginResponseDto, error) {
	access_id, refresh_id := uuid.New().String(), uuid.New().String()
	access_exp, _ := time.ParseDuration(c.config.Jwt.AccessExpire)
	refresh_exp, _ := time.ParseDuration(c.config.Jwt.RefreshExpire)

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &t.AccessClaim{UID: access_id, RID: refresh_id, UserId: user.ID, Exp: time.Now().Add(access_exp).Unix()})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, &t.RefreshClaim{UID: refresh_id, Exp: time.Now().Add(refresh_exp).Unix()})

	access_token, err := access.SignedString([]byte(c.config.Jwt.AccessSecret))
	if err != nil {
		return nil, err
	}

	refresh_token, err := refresh.SignedString([]byte(c.config.Jwt.RefreshSecret))
	if err != nil {
		return nil, err
	}

	defer func() {
		ctx := context.Background()
		c.redis.Set(ctx, access_id, user.UUID, access_exp)
		c.redis.Set(ctx, refresh_id, user.UUID, refresh_exp)
	}()

	return &response.LoginResponseDto{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User:         *common.NewResponse[u.UserAdvancedResponseDto](user),
	}, nil
}

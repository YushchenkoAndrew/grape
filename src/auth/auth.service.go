package auth

import (
	"context"
	"fmt"
	"grape/src/auth/dto/request"
	r "grape/src/auth/dto/response"
	t "grape/src/auth/types"
	"grape/src/common/config"
	"grape/src/common/service"
	"grape/src/user/dto/response"
	e "grape/src/user/entities"
	"grape/src/user/types"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	db     *gorm.DB
	redis  *redis.Client
	config *config.Config
}

func NewAuthService(s *service.CommonService) *authService {
	return &authService{db: s.DB, redis: s.Redis, config: s.Config}
}

func (c *authService) Login(dto *request.LoginDto) (*r.LoginResponseDto, error) {
	var user *e.UserEntity
	if c.db.Joins("Organization").Limit(1).Find(&user, "name = ? AND status = ?", dto.Name, types.Active); user == nil {
		return nil, fmt.Errorf("user '%s' not found", dto.Name)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Pass)); err != nil {
		return nil, fmt.Errorf("user '%s' not found", dto.Name)
	}

	return c.generate(user)
}

func (c *authService) Refresh(dto *request.RefreshDto) (*r.LoginResponseDto, error) {
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

	var user *e.UserEntity
	if c.db.Joins("Organization").Limit(1).Find(&user, "uuid = ? AND status = ?", user_id, types.Active); user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return c.generate(user)
}

func (c *authService) generate(user *e.UserEntity) (*r.LoginResponseDto, error) {
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
		c.redis.Set(ctx, access_id, uuid.New().String(), access_exp)
		c.redis.Set(ctx, refresh_id, user.UUID, refresh_exp)
	}()

	var res response.UserResponseDto
	copier.Copy(&res, &user)

	return &r.LoginResponseDto{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User:         res,
	}, nil
}

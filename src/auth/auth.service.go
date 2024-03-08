package auth

import (
	"context"
	"fmt"
	"grape/src/auth/dto/request"
	"grape/src/auth/dto/response"
	"grape/src/common/config"
	entities "grape/src/common/entities"
	"grape/src/common/service"
	e "grape/src/user/entities"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (c *authService) Login(dto *request.LoginDto) (*response.LoginResponseDto, error) {
	var user *e.UserEntity
	if c.db.Find(&e.UserEntity{Name: dto.User}).First(&user); user == nil {
		return nil, fmt.Errorf("user '%s' not found", dto.User)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Pass)); err != nil {
		return nil, fmt.Errorf("user '%s' not found", dto.User)
	}

	return c.generate(user)
}

func (c *authService) Refresh(dto *request.RefreshDto) (*response.LoginResponseDto, error) {
	token, err := jwt.Parse(dto.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(c.config.Jwt.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("unauthorized token")
	}

	ctx := context.Background()
	refresh_id, _ := claims["uid"].(string)
	user_id, err := c.redis.Get(ctx, refresh_id).Result()

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	c.redis.Del(ctx, refresh_id)

	var user *e.UserEntity
	if c.db.Find(&e.UserEntity{UuidEntity: &entities.UuidEntity{UUID: user_id}}).First(&user); user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return c.generate(user)
}

func (c *authService) generate(user *e.UserEntity) (*response.LoginResponseDto, error) {
	access_id, refresh_id := uuid.New().String(), uuid.New().String()
	access_exp, _ := time.ParseDuration(c.config.Jwt.AccessExpire)
	refresh_exp, _ := time.ParseDuration(c.config.Jwt.RefreshExpire)

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{"uid": access_id, "rid": refresh_id, "user_id": user.ID, "exp": time.Now().Add(access_exp).Unix()})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{"uid": refresh_id, "exp": time.Now().Add(refresh_exp).Unix()})

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

	return &response.LoginResponseDto{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

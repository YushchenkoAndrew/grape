package test

import (
	"bytes"
	"encoding/json"
	"grape/src/auth/dto/request"
	"grape/src/auth/dto/response"
	"grape/src/common/client"
	"grape/src/common/config"
	"grape/src/common/middleware"
	"grape/src/common/service"
	"net/http"
	"net/http/httptest"
	"testing"

	m "grape/src/common/module"
	_ "grape/src/common/validator"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func SetUpRouter(module func(route *gin.RouterGroup, modules *[]m.ModuleT, s *service.CommonService) m.ModuleT) *gin.Engine {
	gin.SetMode(gin.TestMode)
	cfg := config.NewConfig("configs/", "config_test", "yaml")

	service := &service.CommonService{
		DB:     client.ConnPsql(cfg),
		Redis:  client.ConnRedis(cfg),
		Config: cfg,
	}

	r := gin.New()
	r.Use(gin.Recovery())
	rg := r.Group(cfg.Server.Prefix, middleware.GetMiddleware(service).Default())
	module(rg, &[]m.ModuleT{}, service).Init()
	return r
}

func GetToken(t *testing.T, router *gin.Engine, cfg *config.DatabaseConfig) (string, string) {
	body, _ := json.Marshal(request.LoginDto{Username: cfg.User.Name, Password: cfg.User.Password})
	req, _ := http.NewRequest("POST", "/grape/login", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var res response.LoginResponseDto
	json.Unmarshal(w.Body.Bytes(), &res)

	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)
	return res.AccessToken, res.RefreshToken
}

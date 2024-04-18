package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src/auth/dto/request"
	"grape/src/auth/dto/response"
	"grape/src/common/client"
	"grape/src/common/config"
	common "grape/src/common/dto/response"
	"grape/src/common/middleware"
	"grape/src/common/service"
	pr "grape/src/project/dto/response"
	"net/http"
	"net/http/httptest"
	"testing"

	m "grape/src/common/module"
	_ "grape/src/common/validator"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func SetUpRouter(module func(route *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT) (*gin.Engine, *config.Config) {
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
	module(rg, []m.ModuleT{}, service).Init()
	return r, cfg
}

func GetToken(t *testing.T, router *gin.Engine, cfg *config.Config, db *config.DatabaseConfig) (string, string) {
	body, _ := json.Marshal(request.LoginDto{Username: db.User.Name, Password: db.User.Password})
	req, _ := http.NewRequest("POST", cfg.Server.Prefix+"/auth/login", bytes.NewBuffer(body))

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

func GetProject(t *testing.T, router *gin.Engine, cfg *config.Config, token string) *pr.AdminProjectDetailedResponseDto {
	req, _ := http.NewRequest("GET", cfg.Server.Prefix+"/projects", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var res *common.PageResponseDto[[]pr.AdminProjectBasicResponseDto]
	json.Unmarshal(w.Body.Bytes(), &res)

	require.Greater(t, res.Total, 0)
	require.NotEmpty(t, res.Result[0].Id)

	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/admin/projects/%s", cfg.Server.Prefix, res.Result[0].Id), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var project *pr.AdminProjectDetailedResponseDto
	json.Unmarshal(w.Body.Bytes(), &project)
	return project
}

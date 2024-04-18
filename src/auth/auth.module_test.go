package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src/auth"
	"grape/src/auth/dto/request"
	"grape/src/auth/dto/response"
	"grape/src/common/config"
	"grape/src/common/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var router *gin.Engine
var cfg *config.Config
var db *config.DatabaseConfig

func init() {
	db = config.NewDatabaseConfig("configs/", "database", "yaml")
	router, cfg = test.SetUpRouter(auth.NewAuthModule)
}

func TestLogin(t *testing.T) {
	var token string

	tests := []struct {
		name    string
		method  string
		url     string
		auth    func() string
		body    request.LoginDto
		status  int
		handler func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:   "test correct login",
			method: "POST",
			url:    "/auth/login",
			auth:   func() string { return "" },
			body:   request.LoginDto{Username: db.User.Name, Password: db.User.Password},
			status: http.StatusOK,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LoginResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				token = res.AccessToken
				require.NotEmpty(t, res.AccessToken)
				require.NotEmpty(t, res.RefreshToken)
			},
		},
		{
			name:    "test validate token",
			method:  "GET",
			url:     "/auth/ping",
			auth:    func() string { return token },
			status:  http.StatusOK,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:    "test invalid password",
			method:  "POST",
			url:     "/auth/login",
			auth:    func() string { return "" },
			body:    request.LoginDto{Username: db.User.Name, Password: "invalid"},
			status:  http.StatusUnprocessableEntity,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:    "test invalid username",
			method:  "POST",
			url:     "/auth/login",
			auth:    func() string { return "" },
			body:    request.LoginDto{Username: "invalid", Password: "invalid"},
			status:  http.StatusUnprocessableEntity,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:    "test validate logout",
			method:  "POST",
			url:     "/auth/logout",
			auth:    func() string { return token },
			status:  http.StatusOK,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:    "test check if token is no longer valid",
			method:  "GET",
			url:     "/auth/ping",
			auth:    func() string { return token },
			status:  http.StatusUnauthorized,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.method, cfg.Server.Prefix+test.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.auth()))

			router.ServeHTTP(w, req)

			require.Equal(t, test.status, w.Code)
			test.handler(t, w)
		})
	}
}

func TestRefresh(t *testing.T) {
	_, token := test.GetToken(t, router, cfg, db)

	tests := []struct {
		name     string
		body     request.RefreshDto
		expected int
		validate func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:     "test correct refresh",
			body:     request.RefreshDto{RefreshToken: token},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LoginResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.NotEmpty(t, res.AccessToken)
				require.NotEmpty(t, res.RefreshToken)
			},
		},
		{
			name:     "test invalid refresh",
			body:     request.RefreshDto{RefreshToken: token},
			expected: http.StatusUnprocessableEntity,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", cfg.Server.Prefix+"/auth/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, test.expected, w.Code)
			test.validate(t, w)
		})
	}
}

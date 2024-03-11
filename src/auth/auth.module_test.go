package auth_test

import (
	"bytes"
	"encoding/json"
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
var cfg *config.DatabaseConfig

func init() {
	cfg = config.NewDatabaseConfig("configs/", "database", "yaml")
	router = test.SetUpRouter(auth.NewAuthModule)
}

func TestLogin(t *testing.T) {

	tests := []struct {
		name    string
		body    request.LoginDto
		status  int
		handler func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:   "test correct login",
			body:   request.LoginDto{Username: cfg.User.Name, Password: cfg.User.Password},
			status: http.StatusOK,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LoginResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.NotEmpty(t, res.AccessToken)
				require.NotEmpty(t, res.RefreshToken)
			},
		},
		{
			name:    "test invalid password",
			body:    request.LoginDto{Username: cfg.User.Name, Password: "invalid"},
			status:  http.StatusUnprocessableEntity,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:    "test invalid username",
			body:    request.LoginDto{Username: "invalid", Password: "invalid"},
			status:  http.StatusUnprocessableEntity,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/grape/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, test.status, w.Code)
			test.handler(t, w)
		})
	}
}

func TestRefresh(t *testing.T) {
	_, token := test.GetToken(t, router, cfg)

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
			req, _ := http.NewRequest("POST", "/grape/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, test.expected, w.Code)
			test.validate(t, w)
		})
	}
}

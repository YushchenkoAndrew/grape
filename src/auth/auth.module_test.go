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
			body:   request.LoginDto{Name: cfg.User.Name, Pass: cfg.User.Password},
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
			body:    request.LoginDto{Name: cfg.User.Name, Pass: "invalid"},
			status:  http.StatusUnprocessableEntity,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:    "test invalid username",
			body:    request.LoginDto{Name: "invalid", Pass: "invalid"},
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

	var res response.LoginResponseDto
	t.Run("test correct login", func(t *testing.T) {
		body, _ := json.Marshal(request.LoginDto{Name: cfg.User.Name, Pass: cfg.User.Password})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/grape/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), &res)

		require.Equal(t, http.StatusOK, w.Code)
		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.RefreshToken)
	})

	tests := []struct {
		name    string
		body    request.RefreshDto
		status  int
		handler func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:   "test correct refresh",
			body:   request.RefreshDto{RefreshToken: res.RefreshToken},
			status: http.StatusOK,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LoginResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.NotEmpty(t, res.AccessToken)
				require.NotEmpty(t, res.RefreshToken)
			},
		},
		{
			name:    "test invalid refresh",
			body:    request.RefreshDto{RefreshToken: res.RefreshToken},
			status:  http.StatusUnprocessableEntity,
			handler: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/grape/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, test.status, w.Code)
			test.handler(t, w)
		})
	}
}

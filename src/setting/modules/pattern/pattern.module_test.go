package pattern_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src"
	"grape/src/auth"
	"grape/src/common/config"
	common "grape/src/common/dto/response"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/common/test"
	"grape/src/setting"
	"grape/src/setting/modules/pattern/dto/request"
	"grape/src/setting/modules/pattern/dto/response"
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
	router, cfg = test.SetUpRouter(
		func(route *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
			return src.NewIndexModule(route, []m.ModuleT{
				auth.NewAuthModule(route, []m.ModuleT{}, s),
				setting.NewSettingModule(route, []m.ModuleT{}, s),
			}, s)
		},
	)

}

func TestPatternModule(t *testing.T) {
	validate := func(id, token string, body interface{}) {
		require.NotEmpty(t, id)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/admin/settings/patterns/%s", cfg.Server.Prefix, id), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &body)
	}

	var pattern_id string
	token, _ := test.GetToken(t, router, cfg, db)

	tests := []struct {
		name     string
		method   string
		url      func() string
		auth     string
		body     interface{}
		expected int
		validate func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Pattern create",
			method: "POST",
			url:    func() string { return "/admin/settings/patterns" },
			auth:   token,
			body: request.PatternCreateDto{
				Mode:   "fill",
				Colors: 2,
				Width:  10.5,
				Height: 20.5,
				Path:   "<path d='M20 0L0 10v10l20-10zm0 10v10l20 10V20z'/>~<path d='M20-10V0l20 10V0zm0 30L0 30v10l20-10zm0 10v10l20 10V40z'/>",
				Options: &request.PatternCreateOptionDto{
					MaxStroke:   2.5,
					MaxScale:    5,
					MaxSpacingX: 1.5,
					MaxSpacingY: 1.5,
				},
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.PatternAdvancedResponseDto
				validate(res.Id, token, &entity)

				pattern_id = res.Id
				require.Equal(t, 2, entity.Colors)
				require.Equal(t, "fill", entity.Mode)
				require.Equal(t, float32(10.5), entity.Width)
				require.Equal(t, float32(20.5), entity.Height)
				require.Equal(t, "<path d='M20 0L0 10v10l20-10zm0 10v10l20 10V20z'/>~<path d='M20-10V0l20 10V0zm0 30L0 30v10l20-10zm0 10v10l20 10V40z'/>", entity.Path)
			},
		},
		{
			name:   "Pattern update",
			method: "PUT",
			url: func() string {
				return fmt.Sprintf("/admin/settings/patterns/%s", pattern_id)
			},
			auth: token,
			body: request.PatternUpdateDto{
				Mode:   "stroke",
				Colors: 3,
				Width:  15.5,
				Height: 25.5,
				Path:   "<path />",
				Options: &request.PatternUpdateOptionDto{
					MaxStroke:   3.5,
					MaxScale:    6,
					MaxSpacingX: 2.5,
					MaxSpacingY: 2.5,
				},
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var entity response.PatternAdvancedResponseDto
				validate(pattern_id, token, &entity)

				require.Equal(t, 3, entity.Colors)
				require.Equal(t, "stroke", entity.Mode)
				require.Equal(t, float32(15.5), entity.Width)
				require.Equal(t, float32(25.5), entity.Height)
				require.Equal(t, "<path />", entity.Path)
			},
		},
		{
			name:     "Find one pattern",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/admin/settings/patterns/%s", pattern_id) },
			auth:     token,
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Find all patterns, paginated",
			method:   "GET",
			url:      func() string { return "/admin/settings/patterns?page=1&take=10" },
			auth:     token,
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Find all patterns, filtered",
			method:   "GET",
			url:      func() string { return "/admin/settings/patterns?colors=4&mode=fill" },
			auth:     token,
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:   "Delete pattern",
			method: "DELETE",
			url: func() string {
				return fmt.Sprintf("/admin/settings/patterns/%s", pattern_id)
			},
			auth:     token,
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:   "Delete pattern, should return not found",
			method: "DELETE",
			url: func() string {
				return fmt.Sprintf("/admin/settings/patterns/%s", pattern_id)
			},
			auth:     token,
			expected: http.StatusUnprocessableEntity,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.method, cfg.Server.Prefix+test.url(), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.auth))

			router.ServeHTTP(w, req)

			require.Equal(t, test.expected, w.Code)
			test.validate(t, w)

			if test.auth != "" {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(test.method, cfg.Server.Prefix+test.url(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusUnauthorized, w.Code)
			}
		})
	}
}

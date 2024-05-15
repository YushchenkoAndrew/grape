package palette_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src"
	"grape/src/auth"
	"grape/src/common/config"
	common "grape/src/common/dto/response"
	"grape/src/common/service"
	"grape/src/common/test"
	"grape/src/setting"
	"grape/src/setting/modules/palette/dto/request"
	"grape/src/setting/modules/palette/dto/response"
	"net/http"
	"net/http/httptest"
	"testing"

	m "grape/src/common/module"

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

func TestPaletteModule(t *testing.T) {
	validate := func(id, token string, body interface{}) {
		require.NotEmpty(t, id)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/admin/settings/palettes/%s", cfg.Server.Prefix, id), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &body)
	}

	var palette_id string
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
			name:   "Palette create",
			method: "POST",
			url:    func() string { return "/admin/settings/palettes" },
			auth:   token,
			body: request.PaletteCreateDto{
				Colors: []string{"#355070", "#B56576"},
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.PaletteBasicResponseDto
				validate(res.Id, token, &entity)

				palette_id = res.Id
				require.ElementsMatch(t, []string{"#355070", "#B56576"}, entity.Colors)
			},
		},
		{
			name:   "Palette create should through 400",
			method: "POST",
			url:    func() string { return "/admin/settings/palettes" },
			auth:   token,
			body: request.PaletteCreateDto{
				Colors: []string{"355070", "B56576"},
			},
			expected: http.StatusBadRequest,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:   "Palette update",
			method: "PUT",
			url: func() string {
				return fmt.Sprintf("/admin/settings/palettes/%s", palette_id)
			},
			auth: token,
			body: request.PaletteCreateDto{
				Colors: []string{"#355072", "#B56576", "#70C1B3"},
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var entity response.PaletteBasicResponseDto
				validate(palette_id, token, &entity)

				require.ElementsMatch(t, []string{"#355072", "#B56576", "#70C1B3"}, entity.Colors)
			},
		},
		{
			name:     "Find one palette",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/admin/settings/palettes/%s", palette_id) },
			auth:     token,
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Find all palettes, paginated",
			method:   "GET",
			url:      func() string { return "/admin/settings/palettes?page=1&take=10" },
			auth:     token,
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:   "Delete palette",
			method: "DELETE",
			url: func() string {
				return fmt.Sprintf("/admin/settings/palettes/%s", palette_id)
			},
			auth:     token,
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:   "Delete palettes should return not found",
			method: "DELETE",
			url: func() string {
				return fmt.Sprintf("/admin/settings/palettes/%s", palette_id)
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

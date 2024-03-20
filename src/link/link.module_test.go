package link_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src"
	"grape/src/auth"
	"grape/src/common/config"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/common/test"
	"grape/src/link"
	"grape/src/link/dto/request"
	"grape/src/link/dto/response"
	"grape/src/project"
	pr "grape/src/project/dto/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

var router *gin.Engine
var cfg *config.DatabaseConfig

func init() {
	cfg = config.NewDatabaseConfig("configs/", "database", "yaml")
	router = test.SetUpRouter(
		func(route *gin.RouterGroup, modules []m.ModuleT, s *service.CommonService) m.ModuleT {
			return src.NewIndexModule(route, []m.ModuleT{
				auth.NewAuthModule(route, []m.ModuleT{}, s),
				project.NewProjectModule(route, []m.ModuleT{}, s),
				link.NewLinkModule(route, []m.ModuleT{}, s),
			}, s)
		},
	)
}

func TestLinkModule(t *testing.T) {
	var link_id string
	token, _ := test.GetToken(t, router, cfg)
	project := test.GetProject(t, router, token)

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
			name:     "Link create",
			method:   "POST",
			url:      func() string { return "/admin/links" },
			auth:     token,
			body:     request.LinkCreateDto{Name: "test", Link: "http://test/2", LinkableID: project.Id, LinkableType: "projects"},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LinkBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				link_id = res.Id
				require.NotEmpty(t, res.Id)
				require.Equal(t, "test", res.Name)
				require.Equal(t, "http://test/2", res.Link)
			},
		},
		{
			name:     "Link update name",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/links/%s", link_id) },
			auth:     token,
			body:     request.LinkUpdateDto{Name: "test2"},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LinkBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.Equal(t, "test2", res.Name)
				require.Equal(t, "http://test/2", res.Link)
			},
		},
		{
			name:     "Validate that link attached to project",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/admin/projects/%s", project.Id) },
			auth:     token,
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res pr.AdminProjectDetailedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.Len(t, res.Links, len(project.Links)+1)
				require.Contains(t, lo.Map(res.Links, func(item response.LinkBasicResponseDto, _ int) string { return item.Id }), link_id)
			},
		},
		{
			name:     "Link delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/links/%s", link_id) },
			auth:     token,
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Link delete return not found",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/links/%s", link_id) },
			auth:     token,
			expected: http.StatusUnprocessableEntity,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.method, "/grape"+test.url(), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.auth))

			router.ServeHTTP(w, req)

			require.Equal(t, test.expected, w.Code)
			test.validate(t, w)

			if test.auth != "" {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(test.method, "/grape"+test.url(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusUnauthorized, w.Code)
			}
		})
	}
}

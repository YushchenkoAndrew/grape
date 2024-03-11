package project_test

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
	"grape/src/project"
	"grape/src/project/dto/request"
	"grape/src/project/dto/response"
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
	router = test.SetUpRouter(
		func(route *gin.RouterGroup, modules *[]m.ModuleT, s *service.CommonService) m.ModuleT {
			return src.NewIndexModule(route, &[]m.ModuleT{
				auth.NewAuthModule(route, &[]m.ModuleT{}, s),
				project.NewProjectModule(route, &[]m.ModuleT{}, s),
			}, s)
		},
	)
}

func TestProjectModule(t *testing.T) {
	validate := func(id string, body interface{}) {
		require.NotEmpty(t, id)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/grape/projects/%s", id), nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &body)
	}

	token, _ := test.GetToken(t, router, cfg)
	var project_id string

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
			name:     "Project create",
			method:   "POST",
			url:      func() string { return "/admin/projects" },
			auth:     token,
			body:     request.ProjectCreateDto{Name: "TestProject", Description: "Testing project", Type: "html", Footer: "Test footer"},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.ProjectBasicResponseDto
				validate(res.Id, &entity)

				project_id = res.Id
				require.Equal(t, "TestProject", entity.Name)
				require.Equal(t, "Testing project", entity.Description)
				require.Equal(t, "html", entity.Type)
				require.Equal(t, "Test footer", entity.Footer)
			},
		},
		{
			name:     "Project update",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/projects/%s", project_id) },
			auth:     token,
			body:     request.ProjectUpdateDto{Name: "UpdatedProject", Description: "Updated project description", Type: "markdown", Footer: "Updated footer"},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.ProjectBasicResponseDto
				validate(res.Id, &entity)

				require.Equal(t, "UpdatedProject", entity.Name)
				require.Equal(t, "Updated project description", entity.Description)
				require.Equal(t, "markdown", entity.Type)
				require.Equal(t, "Updated footer", entity.Footer)
			},
		},
		{
			name:     "Project find one",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/projects/%s", project_id) },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Project find all without filter",
			method:   "GET",
			url:      func() string { return "/projects" },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Project find with query",
			method:   "GET",
			url:      func() string { return "/projects?query=test&page=1&take=20" },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Project delete",
			method:   "DELETE",
			auth:     token,
			url:      func() string { return fmt.Sprintf("/admin/projects/%s", project_id) },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Project delete return not found",
			method:   "DELETE",
			auth:     token,
			url:      func() string { return fmt.Sprintf("/admin/projects/%s", project_id) },
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

package link_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src"
	"grape/src/auth"
	"grape/src/common/config"
	req "grape/src/common/dto/request"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/common/test"
	"grape/src/link"
	"grape/src/link/dto/request"
	"grape/src/link/dto/response"
	"grape/src/project"
	pr "grape/src/project/dto/response"
	"grape/src/stage"
	st "grape/src/stage/dto/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
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
				stage.NewStageModule(route, []m.ModuleT{}, s),
				project.NewProjectModule(route, []m.ModuleT{}, s),
				link.NewLinkModule(route, []m.ModuleT{}, s),
			}, s)
		},
	)
}

func TestLinkModule(t *testing.T) {
	token, _ := test.GetToken(t, router, cfg, db)
	project := test.GetProject(t, router, cfg, token)
	task := test.GetTask(t, router, cfg, token)

	validate := func(id string, body interface{}) {
		require.NotEmpty(t, id)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/admin/links/%s", cfg.Server.Prefix, id), nil)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}

		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		if body != nil {
			json.Unmarshal(w.Body.Bytes(), &body)
		}
	}

	// var link_id string
	var links []response.LinkAdvancedResponseDto

	tests := []struct {
		name     string
		method   string
		url      func() string
		auth     string
		body     func() interface{}
		expected int
		validate func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Link create",
			method: "POST",
			url:    func() string { return "/admin/links" },
			auth:   token,
			body: func() interface{} {
				return request.LinkCreateDto{Name: "test", Link: "http://test/2", LinkableID: project.Id, LinkableType: "projects"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LinkBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.LinkAdvancedResponseDto
				validate(res.Id, &entity)
				links = append(links, entity)

				require.NotEmpty(t, res.Id)
				require.Equal(t, "test", res.Name)
				require.Equal(t, "http://test/2", res.Link)
			},
		},
		{
			name:   "Link create",
			method: "POST",
			url:    func() string { return "/admin/links" },
			auth:   token,
			body: func() interface{} {
				return request.LinkCreateDto{Name: "test", Link: "http://test/3", LinkableID: project.Id, LinkableType: "projects"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LinkBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.LinkAdvancedResponseDto
				validate(res.Id, &entity)
				links = append(links, entity)

				require.NotEmpty(t, res.Id)
				require.Equal(t, "test", res.Name)
				require.Equal(t, "http://test/3", res.Link)
				require.Greater(t, entity.Order, links[0].Order)
			},
		},
		{
			name:   "Link create",
			method: "POST",
			url:    func() string { return "/admin/links" },
			auth:   token,
			body: func() interface{} {
				return request.LinkCreateDto{Name: "test", Link: "http://test/3", LinkableID: task.Id, LinkableType: "tasks"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LinkBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.LinkAdvancedResponseDto
				validate(res.Id, &entity)
				links = append(links, entity)

				require.NotEmpty(t, res.Id)
				require.Equal(t, "test", res.Name)
				require.Equal(t, "http://test/3", res.Link)
			},
		},
		{
			name:     "Link update order",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/links/%s/order", links[1].Id) },
			auth:     token,
			body:     func() interface{} { return req.OrderUpdateDto{Position: links[0].Order} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LinkBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.LinkAdvancedResponseDto
				validate(res.Id, &entity)

				var entity2 response.LinkAdvancedResponseDto
				validate(links[0].Id, &entity2)

				require.Equal(t, entity.Order, links[0].Order)
				require.Equal(t, entity2.Order, links[1].Order)
			},
		},
		{
			name:     "Link update revert order",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/links/%s/order", links[1].Id) },
			auth:     token,
			body:     func() interface{} { return req.OrderUpdateDto{Position: links[1].Order} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LinkBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.LinkAdvancedResponseDto
				validate(res.Id, &entity)

				var entity2 response.LinkAdvancedResponseDto
				validate(links[0].Id, &entity2)

				require.Equal(t, entity.Order, links[1].Order)
				require.Equal(t, entity2.Order, links[0].Order)
			},
		},
		{
			name:     "Link update name",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/links/%s", links[0].Id) },
			auth:     token,
			body:     func() interface{} { return request.LinkUpdateDto{Name: "test2"} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.LinkBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.Equal(t, "test2", res.Name)
				require.Equal(t, "http://test/2", res.Link)
			},
		},
		{
			name:     "Validate that link is attached to project",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/admin/projects/%s", project.Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res pr.AdminProjectDetailedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.Len(t, res.Links, len(project.Links)+2)
				require.Contains(t, lo.Map(res.Links, func(item response.LinkAdvancedResponseDto, _ int) string { return item.Id }), links[0].Id)
			},
		},
		{
			name:     "Validate that link is attached to task",
			method:   "GET",
			url:      func() string { return "/admin/stages" },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var stages []st.AdminStageBasicResponseDto
				json.Unmarshal(w.Body.Bytes(), &stages)
				require.Greater(t, len(stages), 0)

				_, found := lo.Find(stages, func(stage st.AdminStageBasicResponseDto) bool {
					e, found := lo.Find(stage.Tasks, func(e st.AdminTaskBasicResponseDto) bool { return e.Id == task.Id })
					if !found {
						return false
					}

					require.Len(t, e.Links, len(task.Links)+1)
					require.Contains(t, lo.Map(e.Links, func(item response.LinkAdvancedResponseDto, _ int) string { return item.Id }), links[2].Id)
					return true
				})

				require.Equal(t, found, true)
			},
		},
		{
			name:     "Link delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/links/%s", links[0].Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var entity response.LinkAdvancedResponseDto
				validate(links[1].Id, &entity)

				require.Equal(t, entity.Order, links[0].Order)
			},
		},
		{
			name:     "Link delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/links/%s", links[1].Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Link delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/links/%s", links[2].Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Link delete return not found",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/links/%s", links[0].Id) },
			body:     func() interface{} { return nil },
			auth:     token,
			expected: http.StatusUnprocessableEntity,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body())

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

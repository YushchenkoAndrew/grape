package tag_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src"
	"grape/src/auth"
	"grape/src/common/config"
	"grape/src/common/dto/response"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/common/test"
	"grape/src/project"
	pr "grape/src/project/dto/response"
	"grape/src/stage"
	st "grape/src/stage/dto/response"
	"grape/src/tag"
	"grape/src/tag/dto/request"
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
				tag.NewTagModule(route, []m.ModuleT{}, s),
			}, s)
		},
	)
}

func TestTagModule(t *testing.T) {
	token, _ := test.GetToken(t, router, cfg, db)
	project := test.GetProject(t, router, cfg, token)
	task := test.GetTask(t, router, cfg, token)

	var tags []response.UuidResponseDto

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
			name:   "Tag create",
			method: "POST",
			url:    func() string { return "/admin/tags" },
			auth:   token,
			body: func() interface{} {
				return request.TagCreateDto{Name: "ProjectTag", TaggableID: project.Id, TaggableType: "projects"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var entity response.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &entity)
				tags = append(tags, entity)

				require.NotEmpty(t, entity.Id)
				require.Equal(t, "ProjectTag", entity.Name)
			},
		},
		{
			name:   "Tag create",
			method: "POST",
			url:    func() string { return "/admin/tags" },
			auth:   token,
			body: func() interface{} {
				return request.TagCreateDto{Name: "TaskTag", TaggableID: task.Id, TaggableType: "tasks"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var entity response.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &entity)
				tags = append(tags, entity)

				require.NotEmpty(t, entity.Id)
				require.Equal(t, "TaskTag", entity.Name)
			},
		},
		{
			name:     "Tag update",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/tags/%s", tags[0].Id) },
			auth:     token,
			body:     func() interface{} { return request.TagUpdateDto{Name: "ProjectTagUpdated"} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.Equal(t, "ProjectTagUpdated", res.Name)
			},
		},
		{
			name:     "Validate that tag is attached to project",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/admin/projects/%s", project.Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res pr.AdminProjectDetailedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.Len(t, res.Tags, len(project.Tags)+1)
				require.Contains(t, lo.Map(res.Tags, func(item response.UuidResponseDto, _ int) string { return item.Id }), tags[0].Id)
			},
		},
		{
			name:     "Validate that tag is attached to task",
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

					require.Len(t, e.Tags, len(task.Tags)+1)
					require.Contains(t, lo.Map(e.Tags, func(item response.UuidResponseDto, _ int) string { return item.Id }), tags[1].Id)
					return true
				})

				require.Equal(t, found, true)
			},
		},
		{
			name:     "Tag delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/tags/%s", tags[0].Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Tag delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/tags/%s", tags[1].Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Tag delete return not found",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/tags/%s", tags[0].Id) },
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

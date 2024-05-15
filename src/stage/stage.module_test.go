package stage_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src"
	"grape/src/auth"
	"grape/src/common/config"
	req "grape/src/common/dto/request"
	common "grape/src/common/dto/response"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/common/test"
	"grape/src/stage"
	"grape/src/stage/dto/request"
	"grape/src/stage/dto/response"
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
			}, s)
		},
	)
}

func TestStageModule(t *testing.T) {
	token, _ := test.GetToken(t, router, cfg, db)

	validate := func(id string) response.AdminStageBasicResponseDto {
		require.NotEmpty(t, id)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/admin/stages", cfg.Server.Prefix), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		router.ServeHTTP(w, req)

		var body []response.AdminStageBasicResponseDto
		require.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &body)

		res, found := lo.Find(body, func(e response.AdminStageBasicResponseDto) bool { return e.Id == id })
		require.Equal(t, found, true)

		return res
	}

	var stages []response.AdminStageBasicResponseDto

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
			name:   "Stage create",
			method: "POST",
			url:    func() string { return "/admin/stages" },
			auth:   token,
			body: func() interface{} {
				return request.StageCreateDto{Name: "TestStage"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				entity := validate(res.Id)
				stages = append(stages, entity)
				require.Equal(t, "TestStage", entity.Name)
			},
		},
		{
			name:   "Stage create",
			method: "POST",
			url:    func() string { return "/admin/stages" },
			auth:   token,
			body: func() interface{} {
				return request.StageCreateDto{Name: "TestStage2"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				entity := validate(res.Id)
				stages = append(stages, entity)
				require.Equal(t, "TestStage2", entity.Name)
			},
		},
		{
			name:   "Task create",
			method: "POST",
			url:    func() string { return fmt.Sprintf("/admin/stages/%s/tasks", stages[0].Id) },
			auth:   token,
			body: func() interface{} {
				return request.TaskCreateDto{Name: "TestTask"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				stages[0] = validate(stages[0].Id)
				entity, found := lo.Find(stages[0].Tasks, func(item response.AdminTaskBasicResponseDto) bool { return item.Id == res.Id })

				require.Equal(t, found, true)
				require.Greater(t, entity.Order, 0)
				require.Equal(t, entity.Name, "TestTask")
				require.Equal(t, entity.Status, "active")
				require.Equal(t, entity.Description, "")
			},
		},
		{
			name:   "Task create",
			method: "POST",
			url:    func() string { return fmt.Sprintf("/admin/stages/%s/tasks", stages[0].Id) },
			auth:   token,
			body: func() interface{} {
				return request.TaskCreateDto{Name: "TestTask2"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				stages[0] = validate(stages[0].Id)
				entity, found := lo.Find(stages[0].Tasks, func(item response.AdminTaskBasicResponseDto) bool { return item.Id == res.Id })

				require.Equal(t, found, true)
				require.Greater(t, entity.Order, 0)
				require.Equal(t, entity.Name, "TestTask2")
				require.Equal(t, entity.Status, "active")
				require.Equal(t, entity.Description, "")
			},
		},
		{
			name:     "Stage update order",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/stages/%s/order", stages[1].Id) },
			auth:     token,
			body:     func() interface{} { return req.OrderUpdateDto{Position: stages[0].Order} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				entity := validate(res.Id)
				entity2 := validate(stages[0].Id)

				require.Equal(t, entity.Order, stages[0].Order)
				require.Equal(t, entity2.Order, stages[1].Order)
			},
		},
		{
			name:     "Stage update revert order",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/stages/%s/order", stages[1].Id) },
			auth:     token,
			body:     func() interface{} { return req.OrderUpdateDto{Position: stages[1].Order} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				entity := validate(res.Id)
				entity2 := validate(stages[0].Id)

				require.Equal(t, entity.Order, stages[1].Order)
				require.Equal(t, entity2.Order, stages[0].Order)
			},
		},
		{
			name:   "Stage update",
			method: "PUT",
			url:    func() string { return fmt.Sprintf("/admin/stages/%s", stages[0].Id) },
			auth:   token,
			body: func() interface{} {
				return request.StageUpdateDto{Name: "UpdatedStage"}
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				entity := validate(res.Id)
				require.Equal(t, "UpdatedStage", entity.Name)
			},
		},
		{
			name:   "Task update",
			method: "PUT",
			url:    func() string { return fmt.Sprintf("/admin/stages/%s/tasks/%s", stages[0].Id, stages[0].Tasks[0].Id) },
			auth:   token,
			body: func() interface{} {
				return request.TaskUpdateDto{Name: "UpdatedTask", Description: "TaskDescription"}
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res common.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				stage := validate(stages[0].Id)
				entity, found := lo.Find(stage.Tasks, func(item response.AdminTaskBasicResponseDto) bool { return item.Id == res.Id })

				require.Equal(t, found, true)
				require.Equal(t, entity.Name, "UpdatedTask")
				require.Equal(t, entity.Description, "TaskDescription")
			},
		},
		{
			name:     "Stages get all",
			method:   "GET",
			url:      func() string { return "/stages" },
			body:     func() interface{} { return nil },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Task delete",
			method:   "DELETE",
			auth:     token,
			url:      func() string { return fmt.Sprintf("/admin/stages/%s/tasks/%s", stages[0].Id, stages[0].Tasks[0].Id) },
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Task delete return not found",
			method:   "DELETE",
			auth:     token,
			url:      func() string { return fmt.Sprintf("/admin/stages/%s/tasks/%s", stages[0].Id, stages[0].Tasks[0].Id) },
			body:     func() interface{} { return nil },
			expected: http.StatusUnprocessableEntity,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Stage delete",
			method:   "DELETE",
			auth:     token,
			url:      func() string { return fmt.Sprintf("/admin/stages/%s", stages[0].Id) },
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				entity := validate(stages[1].Id)
				require.Equal(t, entity.Order, stages[0].Order)
			},
		},
		{
			name:     "Stage delete",
			method:   "DELETE",
			auth:     token,
			url:      func() string { return fmt.Sprintf("/admin/stages/%s", stages[1].Id) },
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {

			},
		},
		{
			name:     "Stage delete return not found",
			method:   "DELETE",
			auth:     token,
			url:      func() string { return fmt.Sprintf("/admin/stages/%s", stages[0].Id) },
			body:     func() interface{} { return nil },
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

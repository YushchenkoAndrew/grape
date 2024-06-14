package context_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src"
	"grape/src/auth"
	"grape/src/common/config"
	req "grape/src/common/dto/request"
	res "grape/src/common/dto/response"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/common/test"
	"grape/src/context"
	"grape/src/context/dto/request"
	"grape/src/context/dto/response"
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
				context.NewContextModule(route, []m.ModuleT{}, s),
			}, s)
		},
	)
}

func TestContextModule(t *testing.T) {
	token, _ := test.GetToken(t, router, cfg, db)
	task := test.GetTask(t, router, cfg, token)

	validate := func(id string, body interface{}) {
		require.NotEmpty(t, id)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/admin/contexts/%s", cfg.Server.Prefix, id), nil)
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

	var contexts []response.ContextAdvancedResponseDto

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
			name:   "Context create",
			method: "POST",
			url:    func() string { return "/admin/contexts" },
			auth:   token,
			body: func() interface{} {
				return request.ContextCreateDto{Name: "test", ContextableID: task.Id, ContextableType: "tasks"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.ContextAdvancedResponseDto
				validate(res.Id, &entity)
				contexts = append(contexts, entity)

				require.NotEmpty(t, res.Id)
				require.Equal(t, "test", res.Name)
			},
		},
		{
			name:   "Context create",
			method: "POST",
			url:    func() string { return "/admin/contexts" },
			auth:   token,
			body: func() interface{} {
				return request.ContextCreateDto{Name: "test2", ContextableID: task.Id, ContextableType: "tasks"}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.ContextAdvancedResponseDto
				validate(res.Id, &entity)
				contexts = append(contexts, entity)

				require.NotEmpty(t, res.Id)
				require.Equal(t, "test2", res.Name)
				require.Greater(t, entity.Order, contexts[0].Order)
			},
		},
		{
			name:   "ContextField create",
			method: "POST",
			url:    func() string { return fmt.Sprintf("/admin/contexts/%s/fields", contexts[0].Id) },
			auth:   token,
			body: func() interface{} {
				return request.ContextFieldCreateDto{Name: "ContextFieldCreated", Options: &map[string]interface{}{"test": "FieldCreated"}}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var context response.ContextAdvancedResponseDto
				validate(contexts[0].Id, &context)
				contexts[0] = context

				entity, found := lo.Find(contexts[0].Fields, func(item response.ContextFieldAdvancedResponseDto) bool { return item.Id == res.Id })

				require.Equal(t, found, true)
				require.Greater(t, entity.Order, 0)
				require.Equal(t, entity.Name, "ContextFieldCreated")

				require.Nil(t, entity.Value)
				require.NotNil(t, entity.Options)
				require.Equal(t, (*entity.Options)["test"], "FieldCreated")
			},
		},
		{
			name:   "ContextField create",
			method: "POST",
			url:    func() string { return fmt.Sprintf("/admin/contexts/%s/fields", contexts[0].Id) },
			auth:   token,
			body: func() interface{} {
				return request.ContextFieldCreateDto{Name: "ContextFieldCreated2", Value: &token, Options: &map[string]interface{}{"test": "FieldCreated"}}
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var context response.ContextAdvancedResponseDto
				validate(contexts[0].Id, &context)
				contexts[0] = context

				entity, found := lo.Find(contexts[0].Fields, func(item response.ContextFieldAdvancedResponseDto) bool { return item.Id == res.Id })

				require.Equal(t, found, true)
				require.Greater(t, entity.Order, 0)
				require.Equal(t, entity.Name, "ContextFieldCreated2")
				require.Equal(t, entity.Value, &token)

				require.NotNil(t, entity.Options)
				require.Equal(t, (*entity.Options)["test"], "FieldCreated")
			},
		},
		{
			name:     "Context update order",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/contexts/%s/order", contexts[1].Id) },
			auth:     token,
			body:     func() interface{} { return req.OrderUpdateDto{Position: contexts[0].Order} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.ContextAdvancedResponseDto
				validate(res.Id, &entity)

				var entity2 response.ContextAdvancedResponseDto
				validate(contexts[0].Id, &entity2)

				require.Equal(t, entity.Order, contexts[0].Order)
				require.Equal(t, entity2.Order, contexts[1].Order)
			},
		},
		{
			name:     "Context update revert order",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/contexts/%s/order", contexts[1].Id) },
			auth:     token,
			body:     func() interface{} { return req.OrderUpdateDto{Position: contexts[1].Order} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.ContextAdvancedResponseDto
				validate(res.Id, &entity)

				var entity2 response.ContextAdvancedResponseDto
				validate(contexts[0].Id, &entity2)

				require.Equal(t, entity.Order, contexts[1].Order)
				require.Equal(t, entity2.Order, contexts[0].Order)
			},
		},
		{
			name:   "ContextField update order",
			method: "PUT",
			url: func() string {
				return fmt.Sprintf("/admin/contexts/%s/fields/%s/order", contexts[0].Id, contexts[0].Fields[1].Id)
			},
			auth:     token,
			body:     func() interface{} { return req.OrderUpdateDto{Position: contexts[0].Fields[0].Order} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.ContextAdvancedResponseDto
				validate(contexts[0].Id, &entity)

				require.Equal(t, entity.Fields[0].Order, contexts[0].Fields[1].Order)
				require.Equal(t, entity.Fields[1].Order, contexts[0].Fields[0].Order)
			},
		},
		{
			name:   "ContextField update revert order",
			method: "PUT",
			url: func() string {
				return fmt.Sprintf("/admin/contexts/%s/fields/%s/order", contexts[0].Id, contexts[0].Fields[1].Id)
			},
			auth:     token,
			body:     func() interface{} { return req.OrderUpdateDto{Position: contexts[0].Fields[1].Order} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.ContextAdvancedResponseDto
				validate(contexts[0].Id, &entity)

				require.Equal(t, entity.Fields[0].Order, contexts[0].Fields[0].Order)
				require.Equal(t, entity.Fields[1].Order, contexts[0].Fields[1].Order)
			},
		},
		{
			name:     "Context update name",
			method:   "PUT",
			url:      func() string { return fmt.Sprintf("/admin/contexts/%s", contexts[0].Id) },
			auth:     token,
			body:     func() interface{} { return request.ContextUpdateDto{Name: "testUpdated"} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.Equal(t, "testUpdated", res.Name)
			},
		},
		{
			name:   "ContextField update",
			method: "PUT",
			url: func() string {
				return fmt.Sprintf("/admin/contexts/%s/fields/%s", contexts[0].Id, contexts[0].Fields[0].Id)
			},
			auth: token,
			body: func() interface{} {
				return request.ContextFieldUpdateDto{Name: "ContextFieldUpdated", Value: "true"}
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res res.UuidResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var context response.ContextAdvancedResponseDto
				validate(contexts[0].Id, &context)

				entity, found := lo.Find(context.Fields, func(item response.ContextFieldAdvancedResponseDto) bool { return item.Id == res.Id })

				require.Equal(t, found, true)
				require.Equal(t, entity.Name, "ContextFieldUpdated")

				require.NotNil(t, entity.Value)
				require.Equal(t, *entity.Value, "true")
			},
		},
		{
			name:     "Validate that context is attached to task",
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

					require.Len(t, e.Contexts, len(task.Contexts)+2)
					require.Contains(t, lo.Map(e.Contexts, func(item response.ContextAdvancedResponseDto, _ int) string { return item.Id }), contexts[0].Id)
					require.Contains(t, lo.Map(e.Contexts, func(item response.ContextAdvancedResponseDto, _ int) string { return item.Id }), contexts[1].Id)
					return true
				})

				require.Equal(t, found, true)
			},
		},
		{
			name:   "ContextField delete",
			method: "DELETE",
			url: func() string {
				return fmt.Sprintf("/admin/contexts/%s/fields/%s", contexts[0].Id, contexts[0].Fields[0].Id)
			},
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var entity response.ContextAdvancedResponseDto
				validate(contexts[0].Id, &entity)

				require.Equal(t, entity.Fields[0].Order, contexts[0].Fields[0].Order)
			},
		},
		{
			name:     "Context delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/contexts/%s", contexts[0].Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var entity response.ContextAdvancedResponseDto
				validate(contexts[1].Id, &entity)

				require.Equal(t, entity.Order, contexts[0].Order)
			},
		},
		{
			name:     "Context delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/contexts/%s", contexts[1].Id) },
			auth:     token,
			body:     func() interface{} { return nil },
			expected: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Context delete return not found",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/contexts/%s", contexts[0].Id) },
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

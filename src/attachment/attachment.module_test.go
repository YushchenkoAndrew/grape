package attachment_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src"
	"grape/src/attachment"
	"grape/src/attachment/dto/response"
	"grape/src/auth"
	"grape/src/common/config"
	req "grape/src/common/dto/request"
	m "grape/src/common/module"
	"grape/src/common/service"
	"grape/src/common/test"
	"grape/src/project"
	pr "grape/src/project/dto/response"
	"grape/src/stage"
	st "grape/src/stage/dto/response"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
				attachment.NewAttachmentModule(route, []m.ModuleT{}, s),
			}, s)
		},
	)
}

func TestAttachmentModule(t *testing.T) {
	token, _ := test.GetToken(t, router, cfg, db)
	project := test.GetProject(t, router, cfg, token)
	task := test.GetTask(t, router, cfg, token)

	validate := func(id string, body interface{}) {
		require.NotEmpty(t, id)
		var route string
		if body != nil {
			route = "/admin"
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/attachments/%s", cfg.Server.Prefix+route, id), nil)
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

	var attachments []response.AttachmentAdvancedResponseDto
	var content string

	tests := []struct {
		name     string
		method   string
		url      func() string
		auth     string
		expected int
		body     func() *bytes.Buffer
		validate func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "Attachment create",
			method: "POST",
			url:    func() string { return "/admin/attachments" },
			auth:   token,
			body: func() *bytes.Buffer {
				file, _ := os.Open(config.GetConfigFile())
				defer file.Close()

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				writer.WriteField("path", "/test/")
				writer.WriteField("attachable_id", project.Id)
				writer.WriteField("attachable_type", "projects")

				part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
				io.Copy(part, file)
				writer.Close()

				content = writer.FormDataContentType()
				return body
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.AttachmentAdvancedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.AttachmentAdvancedResponseDto
				validate(res.Id, &entity)
				validate(res.Id, nil)

				attachments = append(attachments, entity)
				require.Equal(t, "/test/", res.Path)
			},
		},
		{
			name:   "Attachment create",
			method: "POST",
			url:    func() string { return "/admin/attachments" },
			auth:   token,
			body: func() *bytes.Buffer {
				file, _ := os.Open(config.GetConfigFile())
				defer file.Close()

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				writer.WriteField("path", "/test2/")
				writer.WriteField("attachable_id", project.Id)
				writer.WriteField("attachable_type", "projects")

				part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
				io.Copy(part, file)
				writer.Close()

				content = writer.FormDataContentType()
				return body
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.AttachmentAdvancedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.AttachmentAdvancedResponseDto
				validate(res.Id, &entity)
				validate(res.Id, nil)

				attachments = append(attachments, entity)
				require.Equal(t, "/test2/", res.Path)
				require.Greater(t, entity.Order, attachments[0].Order)
			},
		},
		{
			name:   "Attachment create",
			method: "POST",
			url:    func() string { return "/admin/attachments" },
			auth:   token,
			body: func() *bytes.Buffer {
				file, _ := os.Open(config.GetConfigFile())
				defer file.Close()

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				writer.WriteField("path", "/test2/")
				writer.WriteField("attachable_id", task.Id)
				writer.WriteField("attachable_type", "tasks")

				part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
				io.Copy(part, file)
				writer.Close()

				content = writer.FormDataContentType()
				return body
			},
			expected: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.AttachmentAdvancedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.AttachmentAdvancedResponseDto
				validate(res.Id, &entity)
				validate(res.Id, nil)

				attachments = append(attachments, entity)
				require.Equal(t, "/test2/", res.Path)
			},
		},
		{
			name:   "Attachment update order",
			method: "PUT",
			url:    func() string { return fmt.Sprintf("/admin/attachments/%s/order", attachments[1].Id) },
			auth:   token,
			body: func() *bytes.Buffer {
				content = "application/json"
				body, _ := json.Marshal(req.OrderUpdateDto{Position: attachments[0].Order})
				return bytes.NewBuffer(body)
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.AttachmentAdvancedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.AttachmentAdvancedResponseDto
				validate(res.Id, &entity)

				var entity2 response.AttachmentAdvancedResponseDto
				validate(attachments[0].Id, &entity2)

				require.Equal(t, entity.Order, attachments[0].Order)
				require.Equal(t, entity2.Order, attachments[1].Order)
			},
		},
		{
			name:   "Attachment update revert order",
			method: "PUT",
			url:    func() string { return fmt.Sprintf("/admin/attachments/%s/order", attachments[1].Id) },
			auth:   token,
			body: func() *bytes.Buffer {
				content = "application/json"
				body, _ := json.Marshal(req.OrderUpdateDto{Position: attachments[1].Order})
				return bytes.NewBuffer(body)
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.AttachmentAdvancedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				var entity response.AttachmentAdvancedResponseDto
				validate(res.Id, &entity)

				var entity2 response.AttachmentAdvancedResponseDto
				validate(attachments[0].Id, &entity2)

				require.Equal(t, entity.Order, attachments[1].Order)
				require.Equal(t, entity2.Order, attachments[0].Order)
			},
		},
		{
			name:   "Attachment update file name",
			method: "PUT",
			url:    func() string { return fmt.Sprintf("/admin/attachments/%s", attachments[0].Id) },
			auth:   token,
			body: func() *bytes.Buffer {
				file, _ := os.Open(config.GetConfigFile())
				defer file.Close()

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("name", "test.txt")
				writer.Close()

				content = writer.FormDataContentType()
				return body
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.AttachmentAdvancedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				validate(res.Id, nil)
				require.Equal(t, "test.txt", res.Name)
			},
		},
		{
			name:   "Replace file with a new one",
			method: "PUT",
			url:    func() string { return fmt.Sprintf("/admin/attachments/%s", attachments[0].Id) },
			auth:   token,
			body: func() *bytes.Buffer {
				file, _ := os.Open(config.GetConfigFile())
				defer file.Close()

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("name", "test")
				writer.Close()

				content = writer.FormDataContentType()
				return body
			},
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res response.AttachmentAdvancedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				validate(res.Id, nil)
			},
		},
		{
			name:     "Validate that attachment is attached to project",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/admin/projects/%s", project.Id) },
			auth:     token,
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var res pr.AdminProjectDetailedResponseDto
				json.Unmarshal(w.Body.Bytes(), &res)

				require.Len(t, res.Attachments, len(project.Attachments)+2)
				require.Contains(t, lo.Map(res.Attachments, func(item response.AttachmentAdvancedResponseDto, _ int) string { return item.Id }), attachments[0].Id)
				require.Contains(t, lo.Map(res.Attachments, func(item response.AttachmentAdvancedResponseDto, _ int) string { return item.Id }), attachments[1].Id)
			},
		},
		{
			name:     "Validate that attachment is attached to task",
			method:   "GET",
			url:      func() string { return "/admin/stages" },
			auth:     token,
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
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

					require.Len(t, e.Attachments, len(task.Attachments)+1)
					require.Contains(t, lo.Map(e.Attachments, func(item response.AttachmentAdvancedResponseDto, _ int) string { return item.Id }), attachments[2].Id)
					return true
				})

				require.Equal(t, found, true)
			},
		},
		{
			name:     "Attachment get",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/attachments/%s", attachments[0].Id) },
			auth:     "",
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Attachment delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/attachments/%s", attachments[0].Id) },
			auth:     token,
			expected: http.StatusNoContent,
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var entity response.AttachmentAdvancedResponseDto
				validate(attachments[1].Id, &entity)

				require.Equal(t, entity.Order, attachments[0].Order)
			},
		},
		{
			name:     "Attachment delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/attachments/%s", attachments[1].Id) },
			auth:     token,
			expected: http.StatusNoContent,
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Attachment delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/attachments/%s", attachments[2].Id) },
			auth:     token,
			expected: http.StatusNoContent,
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Attachment delete return not found",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/attachments/%s", attachments[0].Id) },
			auth:     token,
			expected: http.StatusUnprocessableEntity,
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.method, cfg.Server.Prefix+test.url(), test.body())
			req.Header.Set("Content-Type", content)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.auth))

			router.ServeHTTP(w, req)

			require.Equal(t, test.expected, w.Code)
			test.validate(t, w)

			if test.auth != "" {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(test.method, cfg.Server.Prefix+test.url(), test.body())
				req.Header.Set("Content-Type", content)

				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusUnauthorized, w.Code)
			}
		})
	}
}

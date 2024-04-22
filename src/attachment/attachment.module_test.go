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
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
				project.NewProjectModule(route, []m.ModuleT{}, s),
				attachment.NewAttachmentModule(route, []m.ModuleT{}, s),
			}, s)
		},
	)
}

func TestAttachmentModule(t *testing.T) {
	token, _ := test.GetToken(t, router, cfg, db)

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
				writer.WriteField("attachable_id", test.GetProject(t, router, cfg, token).Id)
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
				writer.WriteField("attachable_id", test.GetProject(t, router, cfg, token).Id)
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

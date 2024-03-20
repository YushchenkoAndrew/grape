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
var cfg *config.DatabaseConfig

func init() {
	cfg = config.NewDatabaseConfig("configs/", "database", "yaml")
	router = test.SetUpRouter(
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
	validate := func(id string) {
		require.NotEmpty(t, id)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/grape/attachments/%s", id), nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	}

	var attachment_id string
	var content string
	token, _ := test.GetToken(t, router, cfg)

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
				writer.WriteField("attachable_id", test.GetProject(t, router, token).Id)
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

				validate(res.Id)
				attachment_id = res.Id
				require.Equal(t, "/test/", res.Path)
			},
		},
		{
			name:   "Attachment update file name",
			method: "PUT",
			url:    func() string { return fmt.Sprintf("/admin/attachments/%s", attachment_id) },
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

				validate(res.Id)
				require.Equal(t, "test.txt", res.Name)
			},
		},
		{
			name:   "Replace file with a new one",
			method: "PUT",
			url:    func() string { return fmt.Sprintf("/admin/attachments/%s", attachment_id) },
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

				validate(res.Id)
			},
		},
		{
			name:     "Attachment get",
			method:   "GET",
			url:      func() string { return fmt.Sprintf("/attachments/%s", attachment_id) },
			auth:     "",
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			expected: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Attachment delete",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/attachments/%s", attachment_id) },
			auth:     token,
			expected: http.StatusNoContent,
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name:     "Attachment delete return not found",
			method:   "DELETE",
			url:      func() string { return fmt.Sprintf("/admin/attachments/%s", attachment_id) },
			auth:     token,
			expected: http.StatusUnprocessableEntity,
			body:     func() *bytes.Buffer { return &bytes.Buffer{} },
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.method, "/grape"+test.url(), test.body())
			req.Header.Set("Content-Type", content)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.auth))

			router.ServeHTTP(w, req)

			require.Equal(t, test.expected, w.Code)
			test.validate(t, w)

			if test.auth != "" {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(test.method, "/grape"+test.url(), test.body())
				req.Header.Set("Content-Type", content)

				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusUnauthorized, w.Code)
			}
		})
	}
}

package void

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grape/src/common/service"
	"grape/src/void/dto/api/response"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

type VoidService struct {
	http func(method, path, content string, body io.Reader) (*http.Response, error)
}

func NewVoidService(s *service.CommonService) *VoidService {
	return &VoidService{
		http: func(method, path, content string, body io.Reader) (*http.Response, error) {
			client := &http.Client{Timeout: 30 * time.Second}

			url, _ := url.Parse(s.Config.Void.Url)
			req, err := http.NewRequest(method, url.JoinPath(path).String(), body)
			if err != nil {
				return nil, err
			}

			req.Header.Set("Content-Type", content)
			req.SetBasicAuth(s.Config.Void.Username, s.Config.Void.Password)
			res, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			return res, nil
		},
	}
}

func (c *VoidService) Get(path string) (string, []byte, error) {
	res, err := c.http("GET", path, "application/json", nil)
	if err != nil {
		return "", nil, err
	}

	data, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf(res.Status)
	}

	return res.Header.Get("Content-Type"), data, err
}

func (c *VoidService) Save(path, filename string, file io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filename)
	io.Copy(part, file)
	writer.Close()

	res, err := c.http("POST", path, writer.FormDataContentType(), body)
	if err != nil {
		return err
	}

	_, err = c.response(res.Body)
	return err
}

func (c *VoidService) Delete(path string) (*response.VoidApiResponseDto, error) {
	res, err := c.http("DELETE", path, "application/json", nil)
	if err != nil {
		return nil, err
	}

	return c.response(res.Body)
}

func (c *VoidService) Rename(src, dst, filename string) error {
	res, err := c.http("GET", src, "application/json", nil)
	if err != nil {
		return err
	}

	if _, err := c.Delete(src); err != nil {
		return err
	}

	return c.Save(dst, filename, res.Body)
}

func (c *VoidService) response(body io.Reader) (*response.VoidApiResponseDto, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	res := &response.VoidApiResponseDto{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if res.Status != "OK" {
		return nil, fmt.Errorf(res.Message)
	}

	return res, nil
}

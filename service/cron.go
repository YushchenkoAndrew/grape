package service

import (
	"api/config"
	"api/helper"
	m "api/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CronService struct {
}

func NewCronService() *CronService {
	return &CronService{}
}

// func (*CronService) query(dto *m.CronQueryDto) string {
// 	var result []string

// 	if len(dto.ID) > 0 {
// 		result = append(result, fmt.Sprintf("id=%s", dto.ID))
// 	}

// 	if len(dto.CronTime) > 0 {
// 		result = append(result, fmt.Sprintf("cron_time=%s", dto.CronTime))
// 	}

// 	if len(dto.URL) > 0 {
// 		result = append(result, fmt.Sprintf("url=%s", dto.URL))
// 	}

// 	if len(dto.Method) > 0 {
// 		result = append(result, fmt.Sprintf("method=%s", dto.Method))
// 	}

// 	if !dto.CreatedFrom.IsZero() {
// 		result = append(result, fmt.Sprintf("created_from=%s", dto.CreatedFrom))
// 	}

// 	if !dto.CreatedTo.IsZero() {
// 		result = append(result, fmt.Sprintf("created_to=%s", dto.CreatedTo))
// 	}

// 	return strings.Join(result, "&")
// }

func (c *CronService) Create(dto *m.CronDto) (*m.CronEntity, error) {
	var body, err = json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	var salt, token = helper.BotToken()
	if req, err = http.NewRequest("POST", fmt.Sprintf("%s/cron/subscribe?key=%s", config.ENV.BotUrl, token), bytes.NewBuffer(body)); err != nil {
		return nil, err
	}

	req.Header.Set("X-Custom-Header", salt)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Bot request error")
	}

	var model m.CronEntity
	if err = json.NewDecoder(res.Body).Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *CronService) Read(query *m.CronQueryDto) ([]m.CronEntity, error) {
	return nil, fmt.Errorf("Not implimented")
}

func (c *CronService) Update(query *m.CronQueryDto, dto *m.CronDto) (*m.CronEntity, error) {
	var body, err = json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	var salt, token = helper.BotToken()
	if req, err = http.NewRequest("PUT", fmt.Sprintf("%s/cron/subscribe?key=%s&id=%s", config.ENV.BotUrl, token, query.ID), bytes.NewBuffer(body)); err != nil {
		return nil, err
	}

	req.Header.Set("X-Custom-Header", salt)
	req.Header.Set("Content-Type", "application/json")

	var res *http.Response

	client := &http.Client{}
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bot request error")
	}

	var model m.CronEntity
	if err = json.NewDecoder(res.Body).Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (c *CronService) Delete(query *m.CronQueryDto) error {
	var salt, token = helper.BotToken()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/cron/subscribe?key=%s&id=%s", config.ENV.BotUrl, token, query.ID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Custom-Header", salt)
	req.Header.Set("Content-Type", "application/json")

	var res *http.Response

	client := &http.Client{}
	if res, err = client.Do(req); err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Bot request error")
	}

	return nil
}

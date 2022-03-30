package models

import "time"

type CronDto struct {
	CronTime string `json:"cron_time" xml:"cron_time" example:"00 00 00 */1 * *"`
	URL      string `json:"url" xml:"url" example:"http://127.0.0.1:8000/ping"`
	Method   string `json:"method" xml:"method" example:"post"`
	Token    string `json:"token" xml:"token" example:"HELLO_WORLD"`
	Data     string `json:"data,omitempty" xml:"data" example:"{'data' : 'Hello world'}"`
}

type CronEntity struct {
	ID        string  `json:"id" xml:"id" example:"d266389ebf09e1e8a95a5b4286b504b2"`
	CreatedAt string  `json:"created_at" xml:"created_at" example:"Mon Jan 31 2022 00:00:00 GMT+0000 (Coordinated Universal Time)"`
	Exec      CronDto `json:"exec" xml:"exec"`
}

type CronQueryDto struct {
	ID       string `form:"id" example:"d266389ebf09e1e8a95a5b4286b504b2"`
	CronTime string `form:"cron_time,omitempty" example:"00 00 00 */1 * *"`
	URL      string `form:"url,omitempty" example:"http://127.0.0.1:8000/ping"`
	Method   string `form:"method,omitempty" example:"post"`

	CreatedTo   time.Time `form:"created_to,omitempty" time_format:"2006-01-02" example:"2021-08-06"`
	CreatedFrom time.Time `form:"created_from,omitempty" time_format:"2006-01-02" example:"2021-08-06"`
}

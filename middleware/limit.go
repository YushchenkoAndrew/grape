package middleware

import (
	"api/config"
	"api/helper"
	"api/logs"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (o *Middleware) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := "IP:" + c.ClientIP()
		ctx := context.Background()

		rate, err := o.client.Get(ctx, ip).Int()
		if err != nil {
			go o.client.Set(ctx, ip, 1, time.Duration(config.ENV.RateTime)*time.Second)
			return
		}

		go o.client.Incr(ctx, ip)
		if rate >= config.ENV.RateLimit {
			helper.ErrHandler(c, http.StatusTooManyRequests, "Toggled Reqest rate limiter")
			go logs.SendLogs(&logs.Message{
				Stat:    "OK",
				Name:    "API",
				Url:     "/api/refresh",
				File:    "/middleware/limit.go",
				Message: "Jeez man calm down, you've had inaff of traffic, I'm blocking you; " + ip,
			})
		}
	}
}

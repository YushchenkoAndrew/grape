package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"grape/src/common/dto/response"
	e "grape/src/user/entities"
)

func (c *Middleware) Default() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ip := "IP:" + c.ClientIP()
		// ctx := context.Background()

		// rate, err := o.client.Get(ctx, ip).Int()
		// if err != nil {
		// 	go o.client.Set(ctx, ip, 1, time.Duration(config.ENV.RateTime)*time.Second)
		// 	return
		// }

		// go o.client.Incr(ctx, ip)
		// if rate >= config.ENV.RateLimit {
		// 	helper.ErrHandler(c, http.StatusTooManyRequests, "Toggled Reqest rate limiter")
		// 	go client.SendLogs(&client.Message{
		// 		Stat:    "OK",
		// 		Name:    "grape",
		// 		Url:     "/api/refresh",
		// 		File:    "/middleware/limit.go",
		// 		Message: "Jeez man calm down, you've had inaff of traffic, I'm blocking you; " + ip,
		// 	})
		// }

		var org []e.OrganizationEntity
		if c.db.Limit(1).Find(&org, "organizations.default"); len(org) == 0 {
			response.ThrowErr(ctx, http.StatusUnprocessableEntity, "organization not found")
			return
		}

		ctx.Set("user", &e.UserEntity{Organization: org[0]})
		ctx.Next()
	}
}

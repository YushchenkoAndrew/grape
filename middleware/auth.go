package middleware

import (
	"api/config"
	"api/helper"
	"api/logs"
	"api/models"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (o *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearToken := strings.Split(c.Request.Header.Get("Authorization"), " ")

		if len(bearToken) != 2 {
			helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/refresh",
				File:    "/controllers/index.go",
				Message: "It's mine first time so please be gentle",
			})
			return
		}

		token, err := jwt.Parse(bearToken[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(("invalid signing method"))
			}

			if _, ok := t.Claims.(jwt.Claims); !ok && !t.Valid {
				return nil, fmt.Errorf(("expired token"))
			}

			return []byte(config.ENV.AccessSecret), nil
		})

		if err != nil {
			helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/refresh",
				File:    "/controllers/index.go",
				Message: "It's mine first time so please be gentle",
			})
			return
		}

		var userUUID string
		var accessUUID string

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			var ok bool
			if accessUUID, ok = claims["access_uuid"].(string); !ok {
				helper.ErrHandler(c, http.StatusUnprocessableEntity, "Invalid token inforamation")
				return
			}

			if userUUID, ok = claims["user_id"].(string); !ok {
				helper.ErrHandler(c, http.StatusUnprocessableEntity, "Invalid token inforamation")
				return
			}
		}

		// Final check with cache
		ctx := context.Background()
		if cacheUUID, err := o.client.Get(ctx, accessUUID).Result(); err != nil || cacheUUID != userUUID {
			helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
			return
		}

		// before request

		c.Next()

		// after request

	}
}

func (o *Middleware) AuthToken(key string, model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearToken := strings.Split(c.Request.Header.Get("Authorization"), " ")

		if len(bearToken) != 2 {
			helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				File:    "/middleware/auth.go",
				Message: "Ohhh nyo your token is inccorrect",
			})
			return
		}

		if err, items := helper.Getcache(o.db.Where("MD5(CONCAT(SPLIT_PART(name, '$', 1), ?, ?)) = SPLIT_PART(name, '$', 2)", config.ENV.Pepper, bearToken[1]), o.client, key, fmt.Sprintf("TOKEN=%s", bearToken[1]), model); err != nil || items == 0 {
			helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				File:    "/middleware/auth.go",
				Message: "Ohhh nyo your token is inccorrect",
			})
			return
		}

		// before request

		c.Next()

		// after request

	}
}

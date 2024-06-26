package middleware

import (
	"context"
	t "grape/src/auth/types"
	"grape/src/common/dto/response"
	"grape/src/common/types"
	e "grape/src/user/entities"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (c *Middleware) Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer, token, _ := strings.Cut(ctx.Request.Header.Get("Authorization"), " ")

		if bearer != "Bearer" || token == "" {
			response.ThrowErr(ctx, http.StatusUnauthorized, "invalid token information")
			return
		}

		var claim t.AccessClaim
		jwt, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (interface{}, error) {
			return []byte(c.config.Jwt.AccessSecret), nil
		})

		if err != nil {
			response.ThrowErr(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		if !jwt.Valid {
			response.ThrowErr(ctx, http.StatusUnauthorized, "invalid token information")
			return
		}

		user_id, err := c.redis.Get(context.Background(), claim.UID).Result()
		if err != nil {
			response.ThrowErr(ctx, http.StatusUnauthorized, "invalid token information")
			return
		}

		var users []e.UserEntity
		if c.db.Joins("Organization").Limit(1).Find(&users, "users.id = ? AND users.uuid = ? AND users.status = ?", claim.UserId, user_id, types.Active); len(users) == 0 {
			response.ThrowErr(ctx, http.StatusUnauthorized, "user not found")
			return
		}

		ctx.Set("user", &users[0])
		ctx.Set("access_claim", &claim)
		ctx.Next()
	}
}

// func (o *Middleware) AuthToken(key string, model interface{}) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		valid := regexp.MustCompile("^[A-Za-z0-9]+$")
// 		bearToken := strings.Split(c.Request.Header.Get("Authorization"), " ")

// 		if len(bearToken) != 2 || !valid.MatchString(bearToken[1]) {
// 			helper.ErrHandler(c, http.StatusUnauthorized, "Invalid token inforamation")
// 			go client.SendLogs(&client.Message{
// 				Stat:    "ERR",
// 				Name:    "grape",
// 				File:    "/middleware/auth.go",
// 				Message: "Ohhh nyo your token is inccorrect",
// 			})
// 			return
// 		}

// 		if err, items := helper.Getcache(o.db.Where(fmt.Sprintf("MD5(CONCAT(SPLIT_PART(token, '$', 1), '%s', '%s')) = SPLIT_PART(token, '$', 2)", config.ENV.Pepper, bearToken[1])), o.client, key, fmt.Sprintf("TOKEN=%s", bearToken[1]), model); err != nil || items == 0 {
// 			helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
// 			go client.SendLogs(&client.Message{
// 				Stat:    "ERR",
// 				Name:    "grape",
// 				File:    "/middleware/auth.go",
// 				Message: "Ohhh nyo your token is inccorrect",
// 			})
// 			return
// 		}

// 		// before request

// 		c.Next()

// 		// after request

// 	}
// }

package middleware

import "github.com/go-redis/redis/v8"

type Middleware struct {
	client *redis.Client
}

var middleware Middleware

func NewMiddleware(client *redis.Client) *Middleware {
	middleware = Middleware{client: client}
	return &middleware
}

func GetMiddleware() *Middleware {
	return &middleware
}

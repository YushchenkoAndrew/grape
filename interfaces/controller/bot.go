package controller

import "github.com/gin-gonic/gin"

type Bot interface {
	Redis(c *gin.Context)
}

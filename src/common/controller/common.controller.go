package controller

import "github.com/gin-gonic/gin"

type DefaultController interface {
	FindOne(c *gin.Context)
	FindAll(c *gin.Context)

	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

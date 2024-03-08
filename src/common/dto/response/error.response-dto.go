package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Status  int    `json:"status" xml:"status" example:"404"`
	Error   string `json:"error" xml:"error" example:"Not Found"`
	Message string `json:"message" xml:"message" example:"Server side error: Something went wrong"`
}

func (c *Error) Value() interface{} {
	return c
}

func ThrowErr(c *gin.Context, status int, message string) {
	Build(c, status, &Error{
		Status:  status,
		Message: message,
		Error:   http.StatusText(status),
	})
}

package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Error struct {
	Status  int    `json:"status" xml:"status" example:"404"`
	Error   string `json:"error" xml:"error" example:"Not Found"`
	Message string `json:"message" xml:"message" example:"Server side error: Something went wrong"`
}

func (c *Error) Value() interface{} {
	return c
}

func ThrowErr(ctx *gin.Context, status int, message string) {
	Build(ctx, status, &Error{
		Status:  status,
		Message: message,
		Error:   http.StatusText(status),
	})
}

func BadRequest(ctx *gin.Context, err error) {
	status := http.StatusBadRequest
	var fields validator.ValidationErrors
	if !errors.As(err, &fields) {
		ThrowErr(ctx, status, err.Error())
		return
	}

	Build(ctx, status, &Error{
		Status:  status,
		Message: prettier(fields[0]),
		Error:   http.StatusText(status),
	})
}

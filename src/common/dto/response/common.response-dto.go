package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UuidResponseDto struct {
	Id   string `copier:"UUID" json:"id" xml:"id" example:"a3c94c88-944d-4636-86d1-c233bdfaf70e"`
	Name string `json:"name" xml:"name" example:"root"`
}

type PageResponseDto[T any] struct {
	Page    int `json:"page" xml:"page" example:"1"`
	PerPage int `json:"per_page" xml:"per_page" example:"30"`
	Total   int `json:"total" xml:"total" example:"30"`
	Result  T   `json:"result" xml:"result"`
}

func Build(ctx *gin.Context, status int, c interface{}) {
	defer ctx.Abort()

	switch ctx.GetHeader("Accept") {
	case "application/xml":
		ctx.XML(status, c)

	default:
		ctx.JSON(status, c)
	}
}

func Handler[T any](ctx *gin.Context, status int, res T, err error) {
	if err != nil {
		ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	Build(ctx, status, res)
}

func NewResponse[Response any, Entity any](entity *Entity) *Response {
	var res Response
	copier.Copy(&res, &entity)
	return &res
}

func NewResponseMany[Response any, Entity any](entities []Entity) []Response {
	res := []Response{}
	copier.Copy(&res, &entities)
	return res
}

package response

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UuidResponseDto struct {
	Id   string `copier:"UUID" json:"id" xml:"id" example:"uuid"`
	Name string `json:"name" xml:"name" example:"root"`
}

type PageResponseDto[T any] struct {
	Page    int `json:"page" xml:"page" example:"1"`
	PerPage int `json:"per_page" xml:"per_page" example:"30"`
	Total   int `json:"total" xml:"total" example:"30"`
	Result  T   `json:"result" xml:"result"`
}

func Build(ctx *gin.Context, stat int, c interface{}) {
	switch ctx.GetHeader("Accept") {
	case "application/xml":
		ctx.XML(stat, c)

	default:
		ctx.JSON(stat, c)
	}

	ctx.Abort()
}

func NewResponse[Response any, Entity any](entity *Entity) Response {
	var res Response
	copier.Copy(&res, entity)
	return res
}

func NewResponseMany[Response any, Entity any](entities []Entity) []Response {
	res := []Response{}
	for _, e := range entities {
		res = append(res, NewResponse[Response](&e))
	}

	return res
}

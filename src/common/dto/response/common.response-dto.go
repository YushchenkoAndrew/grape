package response

import "github.com/gin-gonic/gin"

type UuidResponseDto struct {
	Id   string `copier:"UUID" json:"id" xml:"id" example:"uuid"`
	Name string `json:"name" xml:"name" example:"root"`
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

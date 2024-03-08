package response

import "github.com/gin-gonic/gin"

func Build(ctx *gin.Context, stat int, c interface{}) {
	switch ctx.GetHeader("Accept") {
	case "application/xml":
		ctx.XML(stat, c)

	default:
		ctx.JSON(stat, c)
	}

	ctx.Abort()
}

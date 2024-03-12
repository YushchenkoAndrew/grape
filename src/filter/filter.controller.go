package filter

import (
	"grape/src/common/dto/response"
	"grape/src/filter/dto/request"
	"grape/src/user/entities"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FilterController struct {
	service *FilterService
}

func NewFilterController(s *FilterService) *FilterController {
	return &FilterController{service: s}
}

// @Tags Filter
// @Summary Trace Ip :ip
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param ip path string true "Client IP"
// @Success 200 {object} response.PageResponseDto[[]response.ProjectBasicResponseDto]
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /trace/{ip} [get]
func (c *FilterController) TraceIp(ctx *gin.Context) {
	ip := ctx.Param("ip")
	user, _ := ctx.Get("user")

	if net.ParseIP(ip).To4() == nil {
		response.ThrowErr(ctx, http.StatusBadRequest, "invalid ip address")
		return
	}

	res, err := c.service.TraceIP(request.NewLocationDto(user.(*entities.UserEntity), &request.LocationDto{IP: []string{ip}}))
	response.Handler(ctx, http.StatusOK, res, err)
}

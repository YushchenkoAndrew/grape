package customer

import (
	"grape/src/common/dto/response"
	"grape/src/customer/dto/request"
	"grape/src/user/entities"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	service *CustomerService
}

func NewCustomerController(s *CustomerService) *CustomerController {
	return &CustomerController{service: s}
}

// @Tags Customer
// @Summary Trace Ip :ip
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param ip path string true "Client IP"
// @Success 200 {object} response.LocationResponseDto
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /trace/{ip} [get]
func (c *CustomerController) TraceIp(ctx *gin.Context) {
	ip := ctx.Param("ip")
	user, _ := ctx.Get("user")

	if net.ParseIP(ip).To4() == nil {
		response.ThrowErr(ctx, http.StatusBadRequest, "invalid ip address")
		return
	}

	res, err := c.service.TraceIP(request.NewLocationDto(user.(*entities.UserEntity), &request.LocationDto{IP: []string{ip}}))
	response.Handler(ctx, http.StatusOK, res, err)
}

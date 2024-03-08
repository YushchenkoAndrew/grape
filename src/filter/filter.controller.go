package filter

import (
	"github.com/gin-gonic/gin"
)

type FilterT interface {
	TraceIp(c *gin.Context)
}

type filterController struct {
	service *filterService
}

func NewFilterController(s *filterService) FilterT {
	return &filterController{service: s}
}

// @Summary Trace Ip :ip
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param ip path string true "Client IP"
// @Success 200 {object} m.Success{result=[]m.GeoIpLocations}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /trace/{ip} [get]
func (o *filterController) TraceIp(c *gin.Context) {
	// var ip string
	// if ip = c.Param("ip"); ip == "" {
	// 	helper.CreateErr(c, http.StatusBadRequest, "Incorrect ip value")
	// 	return
	// }

	// models, err := o.service.TraceIP(ip)
	// if err != nil {
	// 	helper.CreateErr(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// helper.ResHandler(c, http.StatusOK, &m.Success{
	// 	Status: "OK",
	// 	Result: models,
	// 	Items:  len(models),
	// })
}

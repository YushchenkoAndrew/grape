package link

import (
	"grape/src/common/dto/response"
	"grape/src/link/dto/request"
	"grape/src/user/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkController struct {
	service *LinkService
}

func NewLinkController(s *LinkService) *LinkController {
	return &LinkController{service: s}
}

func (c *LinkController) dto(ctx *gin.Context, init ...*request.LinkDto) *request.LinkDto {
	user, _ := ctx.Get("user")
	return request.NewLinkDto(user.(*entities.UserEntity), init...)
}

// @Tags Link
// @Summary Create attachment
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model formData request.LinkCreateDto true
// @Success 201 {object} response.LinkBasicResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/links [post]
func (c *LinkController) Create(ctx *gin.Context) {
	var body request.LinkCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Link
// @Summary Update Link
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Link id"
// @Param model formData request.LinkUpdateDto true
// @Success 200 {object} response.LinkBasicResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/links/{id} [put]
func (c *LinkController) Update(ctx *gin.Context) {
	var body request.LinkUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	dto := c.dto(ctx, &request.LinkDto{LinkIds: []string{ctx.Param("id")}})
	res, err := c.service.Update(dto, &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Link
// @Summary Delete attachment
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path string true "Link id"
// @Success 204
// @failure 422 {object} response.Error
// @Router /admin/links/{id} [delete]
func (c *LinkController) Delete(ctx *gin.Context) {
	res, err := c.service.Delete(
		c.dto(ctx, &request.LinkDto{LinkIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

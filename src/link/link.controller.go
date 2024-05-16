package link

import (
	req "grape/src/common/dto/request"
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
// @Summary Find link
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Link id"
// @Success 200 {object} response.LinkAdvancedResponseDto
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/links/{id} [get]
func (c *LinkController) AdminFindOne(ctx *gin.Context) {
	res, err := c.service.AdminFindOne(
		c.dto(ctx, &request.LinkDto{LinkIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Link
// @Summary Create link
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model formData request.LinkCreateDto true "Link data"
// @Success 201 {object} response.LinkBasicResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/links [post]
func (c *LinkController) Create(ctx *gin.Context) {
	var body request.LinkCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Link
// @Summary Update Link
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Link id"
// @Param model body request.LinkUpdateDto true "Link data"
// @Success 200 {object} response.LinkBasicResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/links/{id} [put]
func (c *LinkController) Update(ctx *gin.Context) {
	var body request.LinkUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
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
// @Security BearerAuth
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

// @Tags Link
// @Summary Update Link order position
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Link id"
// @Param model body req.OrderUpdateDto true "Position body"
// @Success 200 {object} response.LinkBasicResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/links/{id}/order [put]
func (c *LinkController) UpdateOrder(ctx *gin.Context) {
	var body req.OrderUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.UpdateOrder(c.dto(ctx, &request.LinkDto{LinkIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

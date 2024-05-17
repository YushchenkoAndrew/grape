package context

import (
	req "grape/src/common/dto/request"
	"grape/src/common/dto/response"
	"grape/src/context/dto/request"
	"grape/src/user/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContextController struct {
	service *ContextService
}

func NewContextController(s *ContextService) *ContextController {
	return &ContextController{service: s}
}

func (c *ContextController) dto(ctx *gin.Context, init ...*request.ContextDto) *request.ContextDto {
	user, _ := ctx.Get("user")
	return request.NewContextDto(user.(*entities.UserEntity), init...)
}

func (c *ContextController) field_dto(ctx *gin.Context, init ...*request.ContextFieldDto) *request.ContextFieldDto {
	user, _ := ctx.Get("user")
	return request.NewContextFieldDto(user.(*entities.UserEntity), init...)
}

// @Tags Context
// @Summary Find context
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Context id"
// @Success 200 {object} response.ContextAdvancedResponseDto
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts/{id} [get]
func (c *ContextController) AdminFindOne(ctx *gin.Context) {
	res, err := c.service.AdminFindOne(
		c.dto(ctx, &request.ContextDto{ContextIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Context
// @Summary Create contexts
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body request.ContextCreateDto true "Context data"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts [post]
func (c *ContextController) Create(ctx *gin.Context) {
	var body request.ContextCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Context
// @Summary Update Context
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Context id"
// @Param model body request.ContextUpdateDto true "Context data"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts/{id} [put]
func (c *ContextController) Update(ctx *gin.Context) {
	var body request.ContextUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	dto := c.dto(ctx, &request.ContextDto{ContextIds: []string{ctx.Param("id")}})
	res, err := c.service.Update(dto, &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Context
// @Summary Delete context
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Context id"
// @Success 204
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts/{id} [delete]
func (c *ContextController) Delete(ctx *gin.Context) {
	res, err := c.service.Delete(
		c.dto(ctx, &request.ContextDto{ContextIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

// @Tags Context
// @Summary Create field at context
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Context id"
// @Param model body request.ContextFieldCreateDto true "Field body"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts/{id}/fields [post]
func (c *ContextController) CreateField(ctx *gin.Context) {
	var body request.ContextFieldCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.CreateField(c.field_dto(ctx, &request.ContextFieldDto{ContextIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Context
// @Summary Update update field at context
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Context id"
// @Param field_id path string true "Field id"
// @Param model body request.ContextFieldUpdateDto true "Field body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts/{id}/fields/{field_id} [put]
func (c *ContextController) UpdateField(ctx *gin.Context) {
	var body request.ContextFieldUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.UpdateField(c.field_dto(ctx, &request.ContextFieldDto{ContextIds: []string{ctx.Param("id")}, ContextFieldIds: []string{ctx.Param("field_id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Context
// @Summary Delete Field in Context
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Context id"
// @Param field_id path string true "Field id"
// @Success 204
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts/{id}/fields/{field_id} [delete]
func (c *ContextController) DeleteField(ctx *gin.Context) {
	res, err := c.service.DeleteField(
		c.field_dto(ctx, &request.ContextFieldDto{ContextIds: []string{ctx.Param("id")}, ContextFieldIds: []string{ctx.Param("field_id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

// @Tags Context
// @Summary Update Context order position
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Context id"
// @Param model body req.OrderUpdateDto true "Position body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts/{id}/order [put]
func (c *ContextController) UpdateOrder(ctx *gin.Context) {
	var body req.OrderUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.UpdateOrder(c.dto(ctx, &request.ContextDto{ContextIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Context
// @Summary Update Field order position
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Context id"
// @Param field_id path string true "Field id"
// @Param model body req.OrderUpdateDto true "Position body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/contexts/{id}/fields/{field_id}/order [put]
func (c *ContextController) UpdateFieldOrder(ctx *gin.Context) {
	var body req.OrderUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.UpdateFieldOrder(c.field_dto(ctx, &request.ContextFieldDto{ContextIds: []string{ctx.Param("id")}, ContextFieldIds: []string{ctx.Param("field_id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

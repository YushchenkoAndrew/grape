package pattern

import (
	"grape/src/common/dto/response"
	"grape/src/setting/modules/pattern/dto/request"
	"grape/src/user/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PatternController struct {
	service *PatternService
}

func NewPatternController(s *PatternService) *PatternController {
	return &PatternController{service: s}
}

func (c *PatternController) dto(ctx *gin.Context, init ...*request.PatternDto) *request.PatternDto {
	user, _ := ctx.Get("user")
	return request.NewPatternDto(user.(*entities.UserEntity), init...)
}

// @Tags Pattern
// @Summary Find all pattern, paginated
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model query request.PatternDto true "Pattern Data"
// @Success 200 {object} response.PageResponseDto[[]response.PatternBasicResponseDto]
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/patterns [get]
func (c *PatternController) FindAll(ctx *gin.Context) {
	dto := c.dto(ctx)

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.FindAll(dto)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Pattern
// @Summary Find one project
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Project id"
// @Success 200 {object} response.PatternBasicResponseDto
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/patterns/{id} [get]
func (c *PatternController) FindOne(ctx *gin.Context) {
	res, err := c.service.FindOne(
		c.dto(ctx, &request.PatternDto{PatternIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Pattern
// @Summary Create Pattern
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body request.PatternCreateDto true "Project body"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/patterns [post]
func (c *PatternController) Create(ctx *gin.Context) {
	var body request.PatternCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Pattern
// @Summary Update Pattern
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Pattern id"
// @Param model body request.PatternUpdateDto true "Pattern body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/patterns/{id} [put]
func (c *PatternController) Update(ctx *gin.Context) {
	var body request.PatternUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Update(c.dto(ctx, &request.PatternDto{PatternIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Pattern
// @Summary Delete Pattern
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Pattern id"
// @Success 204
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/patterns/{id} [delete]
func (c *PatternController) Delete(ctx *gin.Context) {
	res, err := c.service.Delete(
		c.dto(ctx, &request.PatternDto{PatternIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

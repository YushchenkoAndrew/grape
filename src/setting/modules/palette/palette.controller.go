package palette

import (
	"grape/src/common/dto/response"
	"grape/src/setting/modules/palette/dto/request"
	"grape/src/user/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaletteController struct {
	service *PaletteService
}

func NewPaletteController(s *PaletteService) *PaletteController {
	return &PaletteController{service: s}
}

func (c *PaletteController) dto(ctx *gin.Context, init ...*request.PaletteDto) *request.PaletteDto {
	user, _ := ctx.Get("user")
	return request.NewPaletteDto(user.(*entities.UserEntity), init...)
}

// @Tags Palette
// @Summary Find all palettes, paginated
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model query request.PatternDto true "Pattern Data"
// @Success 200 {object} response.PageResponseDto[[]response.PaletteBasicResponseDto]
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/palettes [get]
func (c *PaletteController) FindAll(ctx *gin.Context) {
	dto := c.dto(ctx)

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.FindAll(dto)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Palette
// @Summary Find one palette
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "palette id"
// @Success 200 {object} response.PaletteBasicResponseDto
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/palettes/{id} [get]
func (c *PaletteController) FindOne(ctx *gin.Context) {
	res, err := c.service.FindOne(
		c.dto(ctx, &request.PaletteDto{PaletteIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Palette
// @Summary Create palette
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body request.PaletteCreateDto true "Project body"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/patterns [post]
func (c *PaletteController) Create(ctx *gin.Context) {
	var body request.PaletteCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Palette
// @Summary Update palette
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "palette id"
// @Param model body request.PaletteCreateDto true "palete body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/settings/patterns/{id} [put]
func (c *PaletteController) Update(ctx *gin.Context) {
	var body request.PaletteCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Update(c.dto(ctx, &request.PaletteDto{PaletteIds: []string{ctx.Param("id")}}), &body)
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
func (c *PaletteController) Delete(ctx *gin.Context) {
	res, err := c.service.Delete(
		c.dto(ctx, &request.PaletteDto{PaletteIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

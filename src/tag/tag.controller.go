package tag

import (
	"grape/src/common/dto/response"
	"grape/src/tag/dto/request"
	"grape/src/user/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	service *TagService
}

func NewTagController(s *TagService) *TagController {
	return &TagController{service: s}
}

func (c *TagController) dto(ctx *gin.Context, init ...*request.TagDto) *request.TagDto {
	user, _ := ctx.Get("user")
	return request.NewTagDto(user.(*entities.UserEntity), init...)
}

// @Tags Tag
// @Summary Find tag
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Tag id"
// @Success 200 {object} response.UuidResponseDto
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/tags/{id} [get]
func (c *TagController) AdminFindOne(ctx *gin.Context) {
	res, err := c.service.AdminFindOne(
		c.dto(ctx, &request.TagDto{TagIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Tag
// @Summary Create Tag
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model formData request.TagCreateDto true "Tag data"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/tags [post]
func (c *TagController) Create(ctx *gin.Context) {
	var body request.TagCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Tag
// @Summary Update Tag
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Tag id"
// @Param model body request.TagUpdateDto true "Tag data"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/tags/{id} [put]
func (c *TagController) Update(ctx *gin.Context) {
	var body request.TagUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	dto := c.dto(ctx, &request.TagDto{TagIds: []string{ctx.Param("id")}})
	res, err := c.service.Update(dto, &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Tag
// @Summary Delete Tag
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Tag id"
// @Success 204
// @failure 422 {object} response.Error
// @Router /admin/tags/{id} [delete]
func (c *TagController) Delete(ctx *gin.Context) {
	res, err := c.service.Delete(
		c.dto(ctx, &request.TagDto{TagIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

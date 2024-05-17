package attachment

import (
	"grape/src/attachment/constant"
	"grape/src/attachment/dto/request"
	req "grape/src/common/dto/request"
	"grape/src/common/dto/response"
	"grape/src/user/entities"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AttachmentController struct {
	service *AttachmentService
}

func NewAttachmentController(s *AttachmentService) *AttachmentController {
	return &AttachmentController{service: s}
}

func (c *AttachmentController) dto(ctx *gin.Context, init ...*request.AttachmentDto) *request.AttachmentDto {
	user, _ := ctx.Get("user")
	return request.NewAttachmentDto(user.(*entities.UserEntity), init...)
}

// @Tags Attachment
// @Summary Find attachment context
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path string true "Attachment id"
// @Success 200 {file} file
// @failure 422 {object} response.Error
// @Router /attachments/{id} [get]
func (c *AttachmentController) FindOne(ctx *gin.Context) {
	content, res, err := c.service.FindOne(
		c.dto(ctx, &request.AttachmentDto{AttachmentIds: []string{ctx.Param("id")}}),
	)

	if err != nil {
		response.ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.Data(http.StatusOK, content, res)
}

// @Tags Attachment
// @Summary Find attachment
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Attachment id"
// @Success 200 {object} response.AttachmentAdvancedResponseDto
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/attachments/{id} [get]
func (c *AttachmentController) AdminFindOne(ctx *gin.Context) {
	res, err := c.service.AdminFindOne(
		c.dto(ctx, &request.AttachmentDto{AttachmentIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Attachment
// @Summary Create attachment
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model formData request.AttachmentCreateDto true "File descriptions"
// @Param file formData file true "File"
// @Success 201 {object} response.AttachmentAdvancedResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/attachments/{id} [post]
func (c *AttachmentController) Create(ctx *gin.Context) {
	var body request.AttachmentCreateDto
	if err := ctx.Bind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	file, err := ctx.FormFile(constant.FORM_DATA_FILE_FILE)
	if err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if strings.Contains(file.Filename, " ") {
		response.ThrowErr(ctx, http.StatusBadRequest, "Invalid filename")
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body, file)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Attachment
// @Summary Update attachment
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Attachment id"
// @Param model formData request.AttachmentUpdateDto true "File descriptions"
// @Param file formData file true "File"
// @Success 200 {object} response.AttachmentAdvancedResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/attachments/{id} [put]
func (c *AttachmentController) Update(ctx *gin.Context) {
	var body request.AttachmentUpdateDto
	if err := ctx.Bind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	file, _ := ctx.FormFile(constant.FORM_DATA_FILE_FILE)
	if file != nil && strings.Contains(file.Filename, " ") {
		response.ThrowErr(ctx, http.StatusBadRequest, "Invalid filename")
		return
	}

	dto := c.dto(ctx, &request.AttachmentDto{AttachmentIds: []string{ctx.Param("id")}})
	res, err := c.service.Update(dto, &body, file)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Attachment
// @Summary Delete attachment
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Attachment id"
// @Success 204
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/attachments/{id} [delete]
func (c *AttachmentController) Delete(ctx *gin.Context) {
	res, err := c.service.Delete(
		c.dto(ctx, &request.AttachmentDto{AttachmentIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

// @Tags Attachment
// @Summary Update Attachment order position
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Project id"
// @Param model body req.OrderUpdateDto true "Position body"
// @Success 200 {object} response.AttachmentAdvancedResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/attachments/{id}/order [put]
func (c *AttachmentController) UpdateOrder(ctx *gin.Context) {
	var body req.OrderUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.UpdateOrder(c.dto(ctx, &request.AttachmentDto{AttachmentIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

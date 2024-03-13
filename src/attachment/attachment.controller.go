package attachment

import (
	"grape/src/attachment/constant"
	"grape/src/attachment/dto/request"
	"grape/src/common/dto/response"
	"grape/src/user/entities"
	"net/http"

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
// @Success 200 {object} response.ProjectAdvancedResponseDto
// @failure 422 {object} response.Error
// @Router /attachments/{id} [get]
func (c *AttachmentController) FindOne(ctx *gin.Context) {
	res, err := c.service.FindOne(
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
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/attachments/{id} [post]
func (c *AttachmentController) Create(ctx *gin.Context) {
	var body request.AttachmentCreateDto
	if err := ctx.Bind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	file, err := ctx.FormFile(constant.FORM_DATA_FILE_FILE)
	if err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body, file)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// // @Tags File
// // @Summary Update File by :id
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param id path int true "Instance id"
// // @Param model body m.FileDto true "File Data"
// // @Success 200 {object} m.Success{result=[]m.File}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /file/{id} [put]
func (c *AttachmentController) Update(ctx *gin.Context) {
	// 	var body m.FileDto
	// 	var id = helper.GetID(c, "id")

	// 	if err := c.ShouldBind(&body); err != nil || id == 0 {
	// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
	// 		return
	// 	}

	// 	models, err := o.service.Update(&m.FileQueryDto{ID: uint32(id)}, &m.File{Name: body.Name, Path: body.Path, Type: body.Type, Role: body.Role})
	// 	if err != nil {
	// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
	// 		return
	// 	}

	//	helper.ResHandler(c, http.StatusCreated, &m.Success{
	//		Status: "OK",
	//		Result: models,
	//		Items:  len(models),
	//	})
}

// // @Tags File
// // @Summary Delete File by :id
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param id path int true "Instance id"
// // @Success 200 {object} m.Success{result=[]string{}}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /file/{id} [delete]
func (c *AttachmentController) Delete(ctx *gin.Context) {
	// 	var id = helper.GetID(c, "id")

	// 	if id == 0 {
	// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
	// 		return
	// 	}

	// 	items, err := o.service.Delete(&m.FileQueryDto{ID: uint32(id)})
	// 	if err != nil {
	// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
	// 		return
	// 	}

	//	helper.ResHandler(c, http.StatusOK, &m.Success{
	//		Status: "OK",
	//		Result: []string{},
	//		Items:  items,
	//	})
}

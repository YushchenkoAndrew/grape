package project

import (
	req "grape/src/common/dto/request"
	"grape/src/common/dto/response"
	"grape/src/project/dto/request"
	statistic "grape/src/statistic/dto/request"
	"grape/src/user/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	service *ProjectService
}

func NewProjectController(s *ProjectService) *ProjectController {
	return &ProjectController{service: s}
}

func (c *ProjectController) dto(ctx *gin.Context, init ...*request.ProjectDto) *request.ProjectDto {
	user, _ := ctx.Get("user")
	return request.NewProjectDto(user.(*entities.UserEntity), init...)
}

// @Tags Project
// @Summary Find all project, paginated
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model query request.ProjectDto true "Project Data"
// @Success 200 {object} response.PageResponseDto[[]response.ProjectBasicResponseDto]
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /projects [get]
func (c *ProjectController) FindAll(ctx *gin.Context) {
	dto := c.dto(ctx)

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.FindAll(dto)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Project
// @Summary Find one project
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path string true "Project id"
// @Success 200 {object} response.ProjectBasicResponseDto
// @failure 422 {object} response.Error
// @Router /projects/{id} [get]
func (c *ProjectController) FindOne(ctx *gin.Context) {
	res, err := c.service.FindOne(
		c.dto(ctx, &request.ProjectDto{ProjectIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Project
// @Summary Find all project, paginated
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model query request.ProjectDto true "Project Data"
// @Success 200 {object} response.PageResponseDto[[]response.ProjectAdvancedResponseDto]
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/projects [get]
func (c *ProjectController) AdminFindAll(ctx *gin.Context) {
	dto := c.dto(ctx)

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.AdminFindAll(dto)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Project
// @Summary Find one project
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Project id"
// @Success 200 {object} response.ProjectAdvancedResponseDto
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/projects/{id} [get]
func (c *ProjectController) AdminFindOne(ctx *gin.Context) {
	res, err := c.service.AdminFindOne(
		c.dto(ctx, &request.ProjectDto{ProjectIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Project
// @Summary Create project
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body request.ProjectCreateDto true "Project body"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/projects [post]
func (c *ProjectController) Create(ctx *gin.Context) {
	var body request.ProjectCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Project
// @Summary Update Project
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Project id"
// @Param model body request.ProjectUpdateDto true "Project body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/projects/{id} [put]
func (c *ProjectController) Update(ctx *gin.Context) {
	var body request.ProjectUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Update(c.dto(ctx, &request.ProjectDto{ProjectIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Project
// @Summary Delete Project
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Project id"
// @Success 204
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/projects/{id} [delete]
func (c *ProjectController) Delete(ctx *gin.Context) {
	res, err := c.service.Delete(
		c.dto(ctx, &request.ProjectDto{ProjectIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

// @Tags Project
// @Summary Add 1 to appropriate project statistic
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path string true "Attachment id"
// @Param kind path string true "Attachment id"
// @Param model body request.StatisticUpdateDto true "Statistic update"
// @Success 204
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /projects/{id}/statistics [put]
func (c *ProjectController) UpdateProjectStatistics(ctx *gin.Context) {
	var body statistic.StatisticUpdateDto
	if err := ctx.Bind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	dto := c.dto(ctx, &request.ProjectDto{ProjectIds: []string{ctx.Param("id")}})
	res, err := c.service.UpdateProjectStatistics(dto, &body)
	response.Handler(ctx, http.StatusNoContent, res, err)
}

// @Tags Project
// @Summary Update Project order position
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Project id"
// @Param model body req.OrderUpdateDto true "Position body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/projects/{id}/order [put]
func (c *ProjectController) UpdateOrder(ctx *gin.Context) {
	var body req.OrderUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.UpdateOrder(c.dto(ctx, &request.ProjectDto{ProjectIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

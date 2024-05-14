package stage

import (
	req "grape/src/common/dto/request"
	"grape/src/common/dto/response"
	"grape/src/stage/dto/request"
	"grape/src/user/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StageController struct {
	service *StageService
}

func NewStageController(s *StageService) *StageController {
	return &StageController{service: s}
}

func (c *StageController) dto(ctx *gin.Context, init ...*request.StageDto) *request.StageDto {
	user, _ := ctx.Get("user")
	return request.NewStageDto(user.(*entities.UserEntity), init...)
}

func (c *StageController) task_dto(ctx *gin.Context, init ...*request.TaskDto) *request.TaskDto {
	user, _ := ctx.Get("user")
	return request.NewTaskDto(user.(*entities.UserEntity), init...)
}

// @Tags Stage
// @Summary Display dashboard
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model query request.StageDto true "Stage Data"
// @Success 200 {object} []response.StageBasicResponseDto
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /stages [get]
func (c *StageController) FindAll(ctx *gin.Context) {
	dto := c.dto(ctx)

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.FindAll(dto)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Stage
// @Summary Display admin dashboard
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model query request.StageDto true "Stage Data"
// @Success 200 {object} []response.AdminStageBasicResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages [get]
func (c *StageController) AdminFindAll(ctx *gin.Context) {
	dto := c.dto(ctx)

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.AdminFindAll(dto)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Stage
// @Summary Create stage
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body request.StageCreateDto true "Stage body"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages [post]
func (c *StageController) Create(ctx *gin.Context) {
	var body request.StageCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Create(c.dto(ctx), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Stage
// @Summary Update stage
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Stage id"
// @Param model body request.StageUpdateDto true "Stage body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages/{id} [put]
func (c *StageController) Update(ctx *gin.Context) {
	var body request.StageUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Update(c.dto(ctx, &request.StageDto{StageIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Stage
// @Summary Delete Stage
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Stage id"
// @Success 204
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages/{id} [delete]
func (c *StageController) Delete(ctx *gin.Context) {
	res, err := c.service.Delete(
		c.dto(ctx, &request.StageDto{StageIds: []string{ctx.Param("id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

// @Tags Stage
// @Summary Create task at stage
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Stage id"
// @Param model body request.StageCreateDto true "Stage body"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages/{id}/tasks [post]
func (c *StageController) CreateTask(ctx *gin.Context) {
	var body request.TaskCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.CreateTask(c.task_dto(ctx, &request.TaskDto{StageIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Stage
// @Summary Create update task at stage
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Stage id"
// @Param task_id path string true "Task id"
// @Param model body request.StageCreateDto true "Stage body"
// @Success 201 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages/{id}/tasks/{task_id} [put]
func (c *StageController) UpdateTask(ctx *gin.Context) {
	var body request.TaskUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.UpdateTask(c.task_dto(ctx, &request.TaskDto{StageIds: []string{ctx.Param("id")}, TaskIds: []string{ctx.Param("task_id")}}), &body)
	response.Handler(ctx, http.StatusCreated, res, err)
}

// @Tags Stage
// @Summary Delete Stage
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Stage id"
// @Param task_id path string true "Task id"
// @Success 204
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages/{id}/tasks/{task_id} [delete]
func (c *StageController) DeleteTask(ctx *gin.Context) {
	res, err := c.service.DeleteTask(
		c.task_dto(ctx, &request.TaskDto{StageIds: []string{ctx.Param("id")}, TaskIds: []string{ctx.Param("task_id")}}),
	)

	response.Handler(ctx, http.StatusNoContent, res, err)
}

// @Tags Stage
// @Summary Update Stage order position
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Stage id"
// @Param model body req.OrderUpdateDto true "Position body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages/{id}/order [put]
func (c *StageController) UpdateOrder(ctx *gin.Context) {
	var body req.OrderUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.UpdateOrder(c.dto(ctx, &request.StageDto{StageIds: []string{ctx.Param("id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Stage
// @Summary Update Task order position & move task to another stage
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "Stage id"
// @Param task_id path string true "Task id"
// @Param model body req.OrderUpdateDto true "Position body"
// @Success 200 {object} response.UuidResponseDto
// @failure 400 {object} response.Error
// @failure 401 {object} response.Error
// @failure 422 {object} response.Error
// @Router /admin/stages/{id}/tasks/{task_id}/order [put]
func (c *StageController) UpdateTaskOrder(ctx *gin.Context) {
	var body req.OrderUpdateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.UpdateTaskOrder(c.task_dto(ctx, &request.TaskDto{StageIds: []string{ctx.Param("id")}, TaskIds: []string{ctx.Param("task_id")}}), &body)
	response.Handler(ctx, http.StatusOK, res, err)
}

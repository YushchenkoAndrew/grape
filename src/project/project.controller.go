package project

import (
	c "grape/src/common/controller"
	"grape/src/common/dto/response"
	"grape/src/project/dto/request"
	"grape/src/user/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type projectController struct {
	service *ProjectService
}

func NewProjectController(s *ProjectService) c.CommonController {
	return &projectController{service: s}
}

// @Tags Project
// @Summary Read Project by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Success 200 {object} response.PageResponseDto[[]response.ProjectBasicResponseDto]
// @failure 429 {object} response.Error
// @failure 400 {object} response.Error
// @Router /projects [get]
func (c *projectController) FindAll(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	dto := request.NewProjectDto(user.(*entities.UserEntity))

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.FindAll(dto)
	if err != nil {
		response.ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.Build(ctx, http.StatusOK, res)
}

// @Tags Project
// @Summary Read Project by it's name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path string true "Project id"
// @Success 200 {object} response.ProjectBasicResponseDto
// @failure 422 {object} response.Error
// @Router /projects/{id} [get]
func (c *projectController) FindOne(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	dto := request.NewProjectDto(user.(*entities.UserEntity), &request.ProjectDto{ProjectIds: []string{ctx.Param("id")}})
	// dto.ProjectIds = []string{ctx.Param("id")}

	res, err := c.service.FindOne(dto)
	if err != nil {
		response.ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.Build(ctx, http.StatusOK, res)
}

// @Tags Project
// @Summary Create file by project id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body m.ProjectDto true "Project Data"
// @Success 201 {object} m.Success{result=[]m.Project}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /projects [post]
func (c *projectController) Create(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	dto := request.NewProjectDto(user.(*entities.UserEntity))

	var body request.ProjectCreateDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Create(dto, &body)
	if err != nil {
		response.ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.Build(ctx, http.StatusCreated, res)
}

// @Tags Project
// @Summary Update Project by it's name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Project Name"
// @Param model body m.ProjectDto true "Project without File Data"
// @Success 200 {object} m.Success{result=[]m.Project}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 416 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /project/{name} [put]
func (o *projectController) Update(c *gin.Context) {
	// var body m.ProjectDto
	// var name = c.Param("name")

	// if err := c.ShouldBind(&body); err != nil || name == "" {
	// 	helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { name: %t }", name == ""))
	// 	return
	// }

	// models, err := o.service.Project.Update(&m.ProjectQueryDto{Name: name}, &m.Project{Name: body.Name, Title: body.Title, Flag: body.Flag, Desc: body.Desc, Note: body.Note})
	// if err != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// helper.ResHandler(c, http.StatusCreated, &m.Success{
	// 	Status: "OK",
	// 	Result: models,
	// 	Items:  len(models),
	// })
}

// @Tags Project
// @Summary Delete Project and Files by it's name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Project Name"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 416 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /project/{name} [delete]
func (o *projectController) Delete(c *gin.Context) {
	// var name = c.Param("name")
	// if name == "" {
	// 	helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { name: %t }", false))
	// 	return
	// }

	// var query = m.ProjectQueryDto{Name: name}
	// models, err := o.service.Project.Read(&query)
	// if err != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// for _, item := range models {
	// 	o.service.Link.Delete(&m.LinkQueryDto{ProjectID: item.ID})
	// 	o.service.File.Delete(&m.FileQueryDto{ProjectID: item.ID})
	// 	o.service.Metrics.Delete(&m.MetricsQueryDto{ProjectID: item.ID})

	// 	// NOTE: Delete subscription at two places at the same time
	// 	model, _ := o.service.Subscription.Read(&m.SubscribeQueryDto{ProjectID: item.ID})
	// 	for _, item := range model {
	// 		if err := o.service.Cron.Delete(&m.CronQueryDto{ID: item.CronID}); err == nil {
	// 			o.service.Subscription.Delete(&m.SubscribeQueryDto{ProjectID: item.ID})
	// 		}
	// 	}
	// }

	// items, err := o.service.Project.Delete(&query)
	// if err != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// helper.ResHandler(c, http.StatusOK, &m.Success{
	// 	Status: "OK",
	// 	Result: []string{},
	// 	Items:  items,
	// })
}

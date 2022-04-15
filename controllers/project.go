package controllers

import (
	"api/helper"
	"api/interfaces"
	m "api/models"
	"api/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type projectController struct {
	service *service.FullProjectService
}

func NewProjectController(s *service.FullProjectService) interfaces.Default {
	return &projectController{service: s}
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
// @Router /project [post]
func (o *projectController) CreateOne(c *gin.Context) {
	var body m.ProjectDto
	if err := c.ShouldBind(&body); err != nil || !body.IsOK() {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t }", body.IsOK()))
		return
	}

	var models = m.Project{Name: body.Name, Title: body.Title, Flag: body.Flag, Desc: body.Desc, Note: body.Note}
	if err := o.service.Project.Create(&models); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse Link body. if some field is missing then just skip that
	for _, item := range body.Links {
		if !item.IsOK() {
			continue
		}

		var model = m.Link{ProjectID: models.ID, Name: item.Name, Link: item.Link}
		if err := o.service.Link.Create(&model); err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		models.Links = append(models.Links, model)
	}

	// Parse File body. if some field is missing then just skip that
	for _, item := range body.Files {
		if !item.IsOK() {
			continue
		}

		var model = m.File{ProjectID: models.ID, Name: item.Name, Path: item.Path, Type: item.Type, Role: item.Role}
		if err := o.service.File.Create(&model); err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		models.Files = append(models.Files, model)
	}

	// TODO:
	// Parse Subscription body. if some field is missing then just skip that

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: []m.Project{models},
		Items:  1,
	})
}

// @Tags Project
// @Summary Create Project from list of objects
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body []m.ProjectDto true "List of Project Data"
// @Success 201 {object} m.Success{result=[]m.Project}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /project/list [post]
func (o *projectController) CreateAll(c *gin.Context) {
	var body []m.ProjectDto
	if err := c.ShouldBind(&body); err != nil || len(body) == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t }", len(body) != 0))
		return
	}

	var models = []m.Project{}
	for _, project := range body {
		if !project.IsOK() {
			continue
		}

		var model = m.Project{Name: project.Name, Title: project.Title, Flag: project.Flag, Desc: project.Desc, Note: project.Note}
		if err := o.service.Project.Create(&model); err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Parse Link body. if some field is missing then just skip that
		for _, item := range project.Links {
			if !item.IsOK() {
				var link = m.Link{ProjectID: model.ID, Name: item.Name, Link: item.Link}
				if err := o.service.Link.Create(&link); err != nil {
					helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
					return
				}

				model.Links = append(model.Links, link)
			}
		}

		// Parse File body. if some field is missing then just skip that
		for _, item := range project.Files {
			if item.IsOK() {
				var file = m.File{ProjectID: model.ID, Name: item.Name, Path: item.Path, Type: item.Type, Role: item.Role}
				if err := o.service.File.Create(&file); err != nil {
					helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
					return
				}

				model.Files = append(model.Files, file)
			}
		}

		// TODO:
		// Parse Subscription body. if some field is missing then just skip that

		models = append(models, model)
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Project
// @Summary Read Project by it's name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param name path string true "Project Name"
// @Success 200 {object} m.Success{result=[]m.Project}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /project/{name} [get]
func (o *projectController) ReadOne(c *gin.Context) {
	var name = c.Param("name")
	if name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { name: %t }", false))
		return
	}

	models, err := o.service.Project.Read(&m.ProjectQueryDto{Name: name})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	for i := 0; i < len(models); i++ {
		models[i].Links, _ = o.service.Link.Read(&m.LinkQueryDto{ProjectID: models[i].ID})
		models[i].Files, _ = o.service.File.Read(&m.FileQueryDto{ProjectID: models[i].ID})
		models[i].Subscription, _ = o.service.Subscription.Read(&m.SubscribeQueryDto{ProjectID: models[i].ID})
		models[i].Metrics, _ = o.service.Metrics.Read(&m.MetricsQueryDto{ProjectID: models[i].ID})
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Project
// @Summary Read Project by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id query int false "Type: '1'"
// @Param name query string false "Name: 'CodeRain'"
// @Param flag query string false "Flag: 'js'"
// @Param created_from query string false "CreatedAt date >= start"
// @Param created_to query string false "CreatedAt date <= end"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Param link[id] query int false "Type: '1'"
// @Param link[name] query string false "Type: 'Name: 'main'"
// @Param link[page] query int false "Page: '0'"
// @Param link[limit] query int false "Limit: '1'"
// @Param file[id] query int false "Type: '1'"
// @Param file[name] query string false "Type: 'index.js'"
// @Param file[role] query string false "Role: 'src'"
// @Param file[path] query string false "Path: '/test'"
// @Param file[page] query int false "Page: '0'"
// @Param file[limit] query int false "Limit: '1'"
// @Success 200 {object} m.Success{result=[]m.Project}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /project [get]
func (o *projectController) ReadAll(c *gin.Context) {
	var query = m.ProjectQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Project.Read(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	for i := 0; i < len(models); i++ {
		models[i].Links, _ = o.service.Link.Read(&m.LinkQueryDto{ID: query.Link.ID, Name: query.Link.Name, ProjectID: models[i].ID, Page: query.Link.Page, Limit: query.Link.Limit})
		models[i].Files, _ = o.service.File.Read(&m.FileQueryDto{ID: query.File.ID, Name: query.File.Name, Role: query.File.Role, Path: query.File.Path, ProjectID: models[i].ID, Page: query.File.Page, Limit: query.File.Limit})
		models[i].Subscription, _ = o.service.Subscription.Read(&m.SubscribeQueryDto{ID: query.Subscription.ID, Name: query.Subscription.Name, CronID: query.Subscription.CronID, ProjectID: models[i].ID, Page: query.Subscription.Page, Limit: query.Subscription.Limit})
		models[i].Metrics, _ = o.service.Metrics.Read(&m.MetricsQueryDto{ID: query.Metrics.ID, Name: query.Metrics.Name, Namespace: query.Metrics.Namespace, ContainerName: query.Metrics.ContainerName, CreatedTo: query.Metrics.CreatedTo, CreatedFrom: query.Metrics.CreatedFrom, Page: query.Metrics.Page, Limit: query.Metrics.Limit})
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: models,
		Page:   query.Page,
		Limit:  query.Limit,
		Items:  len(models),
	})
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
func (o *projectController) UpdateOne(c *gin.Context) {
	var body m.ProjectDto
	var name = c.Param("name")

	if err := c.ShouldBind(&body); err != nil || name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { name: %t }", name == ""))
		return
	}

	models, err := o.service.Project.Update(&m.ProjectQueryDto{Name: name}, &m.Project{Name: body.Name, Title: body.Title, Flag: body.Flag, Desc: body.Desc, Note: body.Note})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Project
// @Summary Update Project by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Name: 'CodeRain'"
// @Param flag query string false "Flag: 'js'"
// @Param created_from query string false "CreatedAt date >= start"
// @Param created_to query string false "CreatedAt date <= end"
// @Param model body m.ProjectDto true "Project without File Data"
// @Success 200 {object} m.Success{result=[]m.File}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /project [put]
func (o *projectController) UpdateAll(c *gin.Context) {
	var query = m.ProjectQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	var body m.ProjectDto
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Project.Update(&query, &m.Project{Name: body.Name, Title: body.Title, Flag: body.Flag, Desc: body.Desc, Note: body.Note})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
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
func (o *projectController) DeleteOne(c *gin.Context) {
	var name = c.Param("name")
	if name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { name: %t }", false))
		return
	}

	var query = m.ProjectQueryDto{Name: name}
	models, err := o.service.Project.Read(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, item := range models {
		o.service.Link.Delete(&m.LinkQueryDto{ProjectID: item.ID})
		o.service.File.Delete(&m.FileQueryDto{ProjectID: item.ID})
		o.service.Metrics.Delete(&m.MetricsQueryDto{ProjectID: item.ID})

		// NOTE: Delete subscription at two places at the same time
		model, _ := o.service.Subscription.Read(&m.SubscribeQueryDto{ProjectID: item.ID})
		for _, item := range model {
			if err := o.service.Cron.Delete(&m.CronQueryDto{ID: item.CronID}); err == nil {
				o.service.Subscription.Delete(&m.SubscribeQueryDto{ProjectID: item.ID})
			}
		}
	}

	items, err := o.service.Project.Delete(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []string{},
		Items:  items,
	})
}

// @Tags Project
// @Summary Delete Project by Query and Files with the same project_id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Name: 'CodeRain'"
// @Param flag query string false "Flag: 'js'"
// @Param created_from query string false "CreatedAt date >= start"
// @Param created_to query string false "CreatedAt date <= end"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 416 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /project [delete]
func (o *projectController) DeleteAll(c *gin.Context) {
	var query = m.ProjectQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Project.Read(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, item := range models {
		o.service.Link.Delete(&m.LinkQueryDto{ProjectID: item.ID, Page: -1})
		o.service.File.Delete(&m.FileQueryDto{ProjectID: item.ID, Page: -1})
		o.service.Metrics.Delete(&m.MetricsQueryDto{ProjectID: item.ID, Page: -1})

		// NOTE: Delete subscription at two places at the same time
		model, _ := o.service.Subscription.Read(&m.SubscribeQueryDto{ProjectID: item.ID, Page: -1})
		for _, item := range model {
			if err := o.service.Cron.Delete(&m.CronQueryDto{ID: item.CronID}); err == nil {
				o.service.Subscription.Delete(&m.SubscribeQueryDto{ProjectID: item.ID, Page: -1})
			}
		}
	}

	items, err := o.service.Project.Delete(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []string{},
		Items:  items,
	})
}

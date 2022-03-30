package controllers

import (
	"api/config"
	"api/helper"
	"api/interfaces"
	m "api/models"
	"api/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type linkController struct {
	service *service.LinkService
}

func NewLinkController(s *service.LinkService) interfaces.Default {
	return &linkController{service: s}
}

// @Tags Link
// @Summary Create link by project id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param model body m.LinkDto true "Link info"
// @Success 201 {object} m.Success{result=[]m.Link}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /link/{id} [post]
func (o *linkController) CreateOne(c *gin.Context) {
	var body m.LinkDto
	var id = helper.GetID(c)

	if err := c.ShouldBind(&body); err != nil || !body.IsOK() || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, id: %t }", body.IsOK(), id != 0))
		return
	}

	var model = m.Link{ProjectID: uint32(id), Name: body.Name, Link: body.Link}
	if err := o.service.Create(&model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: []m.Link{model},
		Items:  1,
	})
}

// @Tags Link
// @Summary Create File from list of objects
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project id"
// @Param model body []m.LinkDto true "List of Links info"
// @Success 201 {object} m.Success{result=[]m.Link}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /link/list/{id} [post]
func (o *linkController) CreateAll(c *gin.Context) {
	var body []m.LinkDto
	var id = helper.GetID(c)

	if err := c.ShouldBind(&body); err != nil || len(body) == 0 || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, id: %t }", len(body) != 0, id != 0))
		return
	}

	var models = []m.Link{}
	for _, item := range body {
		if item.IsOK() {
			var model = m.Link{ProjectID: uint32(id), Name: item.Name, Link: item.Link}
			if err := o.service.Create(&model); err != nil {
				helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
				return
			}

			models = append(models, model)
		}
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Link
// @Summary Read Link by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path int true "Instance id"
// @Success 200 {object} m.Success{result=[]m.Link}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /link/{id} [get]
func (o *linkController) ReadOne(c *gin.Context) {
	var id = helper.GetID(c)

	if id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	models, err := o.service.Read(&m.LinkQueryDto{ID: uint32(id)})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: models,
		Items:  len(models),
	})
}

// @Tags Link
// @Summary Read Link by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param project_id query string false "ProjectID: '1'"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} m.Success{result=[]m.Link}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /link [get]
func (o *linkController) ReadAll(c *gin.Context) {
	var query = m.LinkQueryDto{Page: 0, Limit: config.ENV.Limit}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Read(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: models,
		Page:   query.Page,
		Limit:  query.Limit,
		Items:  len(models),
	})
}

// @Tags Link
// @Summary Update Link by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body m.LinkDto true "Link Data"
// @Success 200 {object} m.Success{result=[]m.Link}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /link/{id} [put]
func (o *linkController) UpdateOne(c *gin.Context) {
	var body m.LinkDto
	var id = helper.GetID(c)

	if err := c.ShouldBind(&body); err != nil || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	models, err := o.service.Update(&m.LinkQueryDto{ID: uint32(id)}, &m.Link{Name: body.Name, Link: body.Link, ID: uint32(id)})
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

// @Tags Link
// @Summary Update Link by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param project_id query string false "ProjectID: '1'"
// @Param model body m.LinkDto true "Link Data"
// @Success 200 {object} m.Success{result=[]m.Link}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /link [put]
func (o *linkController) UpdateAll(c *gin.Context) {
	var query = m.LinkQueryDto{}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	var body m.LinkDto
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Update(&query, &m.Link{Name: body.Name, Link: body.Link})
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

// @Tags Link
// @Summary Delete Link by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /link/{id} [delete]
func (o *linkController) DeleteOne(c *gin.Context) {
	var id = helper.GetID(c)

	if id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	items, err := o.service.Delete(&m.LinkQueryDto{ID: uint32(id)})
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

// @Tags Link
// @Summary Delete Link by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param name query string false "Type: 'Name: 'main'"
// @Param project_id query string false "ProjectID: '1'"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /link [delete]
func (o *linkController) DeleteAll(c *gin.Context) {
	var query = m.LinkQueryDto{}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	items, err := o.service.Delete(&query)
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

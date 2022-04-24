package controllers

import (
	"api/helper"
	"api/interfaces/controller"
	"api/interfaces/service"
	m "api/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type fileController struct {
	service service.Default[m.File, m.FileQueryDto]
}

func NewFileController(s service.Default[m.File, m.FileQueryDto]) controller.Default {
	return &fileController{service: s}
}

// @Tags File
// @Summary Create file by project id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param model body m.FileDto true "File Data"
// @Success 201 {object} m.Success{result=[]m.File}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /file/{id} [post]
func (o *fileController) CreateOne(c *gin.Context) {
	var body m.FileDto
	var id = helper.GetID(c, "id")

	if err := c.ShouldBind(&body); err != nil || !body.IsOK() || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, id: %t }", body.IsOK(), id != 0))
		return
	}

	var model = m.File{ProjectID: uint32(id), Name: body.Name, Path: body.Path, Type: body.Type, Role: body.Role}
	if err := o.service.Create(&model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: []m.File{model},
		Items:  1,
	})
}

// @Tags File
// @Summary Create File from list of objects
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project id"
// @Param model body []m.FileDto true "List of File Data"
// @Success 201 {object} m.Success{result=[]m.File}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /file/list/{id} [post]
func (o *fileController) CreateAll(c *gin.Context) {
	var body []m.FileDto
	var id = helper.GetID(c, "id")

	if err := c.ShouldBind(&body); err != nil || len(body) == 0 || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, id: %t }", len(body) != 0, id != 0))
		return
	}

	var models = []m.File{}
	for _, item := range body {
		if item.IsOK() {
			var model = m.File{ProjectID: uint32(id), Name: item.Name, Path: item.Path, Type: item.Type, Role: item.Role}
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

// @Tags File
// @Summary Read File by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path int true "Instance id"
// @Success 200 {object} m.Success{result=[]m.File}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /file/{id} [get]
func (o *fileController) ReadOne(c *gin.Context) {
	var id = helper.GetID(c, "id")

	if id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	models, err := o.service.Read(&m.FileQueryDto{ID: uint32(id)})
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

// @Tags File
// @Summary Read File by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'index.js'"
// @Param role query string false "Role: 'src'"
// @Param path query string false "Path: '/test'"
// @Param project_id query string false "ProjectID: '1'"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} m.Success{result=[]m.File}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /file [get]
func (o *fileController) ReadAll(c *gin.Context) {
	var query = m.FileQueryDto{Page: -1}
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

// @Tags File
// @Summary Update File by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body m.FileDto true "File Data"
// @Success 200 {object} m.Success{result=[]m.File}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /file/{id} [put]
func (o *fileController) UpdateOne(c *gin.Context) {
	var body m.FileDto
	var id = helper.GetID(c, "id")

	if err := c.ShouldBind(&body); err != nil || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	models, err := o.service.Update(&m.FileQueryDto{ID: uint32(id)}, &m.File{Name: body.Name, Path: body.Path, Type: body.Type, Role: body.Role})
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

// @Tags File
// @Summary Update File by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param role query string false "Role: 'src'"
// @Param path query string false "Path: '/test'"
// @Param project_id query string false "ProjectID: '1'"
// @Param model body m.FileDto true "File Data"
// @Success 200 {object} m.Success{result=[]m.File}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /file [put]
func (o *fileController) UpdateAll(c *gin.Context) {
	var query = m.FileQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	var body m.FileDto
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Update(&query, &m.File{Name: body.Name, Path: body.Path, Type: body.Type, Role: body.Role})
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

// @Tags File
// @Summary Delete File by :id
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
// @Router /file/{id} [delete]
func (o *fileController) DeleteOne(c *gin.Context) {
	var id = helper.GetID(c, "id")

	if id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
		return
	}

	items, err := o.service.Delete(&m.FileQueryDto{ID: uint32(id)})
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

// @Tags File
// @Summary Delete File by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param role query string false "Role: 'src'"
// @Param path query string false "Path: '/test'"
// @Param project_id query string false "ProjectID: '1'"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /file [delete]
func (o *fileController) DeleteAll(c *gin.Context) {
	var query = m.FileQueryDto{Page: -1}
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

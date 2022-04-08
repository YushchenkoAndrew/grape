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

type subscriptionController struct {
	service *service.FullSubscriptionService
}

func NewSubscriptionController(s *service.FullSubscriptionService) interfaces.Default {
	return &subscriptionController{service: s}
}

// @Tags Subscription
// @Summary Create Subscription to run operation
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param _ query string false "For more info about query see request: 'GET /operations'"
// @Param model body m.SubscribeDto true "Small info about subscription for k3s"
// @Success 201 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/{id} [post]
func (o *subscriptionController) CreateOne(c *gin.Context) {
	var body m.SubscribeDto
	var id = helper.GetID(c)
	var handler config.Handler

	if err := c.ShouldBind(&body); err != nil || !body.IsOK() || id == 0 || !config.GetOperation(body.Name, &handler) {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, id: %t, operation: %t }", body.IsOK(), id != 0, !config.GetOperation(body.Name, &handler)))
		return
	}

	path, err := helper.FormPathFromHandler(c, &handler)
	if err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	entity, err := o.service.Cron.Create(&m.CronDto{CronTime: body.CronTime, URL: config.ENV.URL + path, Method: handler.Method, Token: helper.HashSecret(helper.GetToken())})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	var model = m.Subscription{ProjectID: uint32(id), Name: body.Name, CronID: entity.ID, CronTime: entity.Exec.CronTime, Method: entity.Exec.Method, Path: path, Token: entity.Exec.Token}
	if err := o.service.Subscription.Create(&model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: []m.Subscription{model},
		Items:  1,
	})
}

// @Tags Subscription
// @Summary Create list of Subscriptions to run operation
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param _ query string false "For more info about query see request: 'GET /operations'"
// @Param model body []m.SubscribeDto true "Small info about subscription for k3s"
// @Success 201 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/list/{id} [post]
func (o *subscriptionController) CreateAll(c *gin.Context) {
	var body []m.SubscribeDto
	var id = helper.GetID(c)
	var handler config.Handler

	if err := c.ShouldBind(&body); err != nil || len(body) == 0 || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, id: %t }", len(body) == 0, id != 0))
		return
	}

	var models []m.Subscription

	for _, item := range body {
		if !item.IsOK() || !config.GetOperation(item.Name, &handler) {
			continue
		}

		if path, err := helper.FormPathFromHandler(c, &handler); err == nil {
			entity, err := o.service.Cron.Create(&m.CronDto{CronTime: item.CronTime, URL: config.ENV.URL + path, Method: handler.Method, Token: helper.HashSecret(helper.GetToken())})
			if err != nil {
				helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
				return
			}

			var model = m.Subscription{ProjectID: uint32(id), Name: item.Name, CronID: entity.ID, CronTime: entity.Exec.CronTime, Method: entity.Exec.Method, Path: path, Token: entity.Exec.Token}
			if err := o.service.Subscription.Create(&model); err != nil {
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

// @Tags Subscription
// @Summary Read subscription by cron_id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "CronID"
// @Success 200 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/{id} [get]
func (o *subscriptionController) ReadOne(c *gin.Context) {
	var id = c.Param("id")
	if id == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { name: %t }", false))
		return
	}

	models, err := o.service.Subscription.Read(&m.SubscribeQueryDto{CronID: id})
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

// @Tags Subscription
// @Summary Read subscription by query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param cron_id query int false "Type: '1'"
// @Param project_id query string false "ProjectID: '1'"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/{id} [get]
func (o *subscriptionController) ReadAll(c *gin.Context) {
	var query = m.SubscribeQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	models, err := o.service.Subscription.Read(&query)
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

// @Tags Subscription
// @Summary Update Subscription by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body m.SubscribeDto true "Small info about subscription for k3s"
// @Success 200 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/{id} [put]
func (o *subscriptionController) UpdateOne(c *gin.Context) {
	var body m.SubscribeDto
	var id = c.Param("id")
	var handler config.Handler

	if err := c.ShouldBind(&body); err != nil || id == "" || !config.GetOperation(body.Name, &handler) {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, id: %t, operation: %t }", body.IsOK(), id == "", !config.GetOperation(body.Name, &handler)))
		return
	}

	path, err := helper.FormPathFromHandler(c, &handler)
	if err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	entity, err := o.service.Cron.Update(&m.CronQueryDto{ID: id}, &m.CronDto{CronTime: body.CronTime, URL: config.ENV.URL + path, Method: handler.Method, Token: helper.HashSecret(helper.GetToken())})
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	models, err := o.service.Subscription.Update(&m.SubscribeQueryDto{CronID: id}, &m.Subscription{Name: body.Name, CronID: entity.ID, CronTime: entity.Exec.CronTime, Method: entity.Exec.Method, Path: path, Token: entity.Exec.Token})
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

// @Tags Subscription
// @Summary Update Subscription by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param cron_id query int false "Type: '1'"
// @Param project_id query string false "ProjectID: '1'"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Param model body m.SubscribeDto true "Small info about subscription for k3s"
// @Success 200 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription [put]
func (o *subscriptionController) UpdateAll(c *gin.Context) {
	var query = m.SubscribeQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	var body m.SubscribeDto
	var handler config.Handler
	if err := c.ShouldBind(&body); err != nil || !config.GetOperation(body.Name, &handler) {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, operation: %t }", body.IsOK(), !config.GetOperation(body.Name, &handler)))
		return
	}

	path, err := helper.FormPathFromHandler(c, &handler)
	if err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	model, err := o.service.Subscription.Read(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	var models = []m.Subscription{}
	for _, item := range model {
		entity, err := o.service.Cron.Update(&m.CronQueryDto{ID: item.CronID}, &m.CronDto{CronTime: body.CronTime, URL: config.ENV.URL + path, Method: handler.Method, Token: helper.HashSecret(helper.GetToken())})
		if err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		model, err := o.service.Subscription.Update(&m.SubscribeQueryDto{CronID: item.CronID}, &m.Subscription{Name: body.Name, CronID: entity.ID, CronTime: entity.Exec.CronTime, Method: entity.Exec.Method, Path: path, Token: entity.Exec.Token})
		if err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		models = append(models, model...)
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: models,
		Page:   query.Page,
		Limit:  query.Limit,
		Items:  len(models),
	})
}

// @Tags Subscription
// @Summary Delete subscription by cron_id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "CronID"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/{id} [delete]
func (o *subscriptionController) DeleteOne(c *gin.Context) {
	var id = c.Param("id")
	if id == "" {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", false))
		return
	}

	if err := o.service.Cron.Delete(&m.CronQueryDto{ID: id}); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	items, err := o.service.Subscription.Delete(&m.SubscribeQueryDto{CronID: id})
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

// @Tags Subscription
// @Summary Delete subscription by query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "CronID"
// @Param name query string false "Type: 'Name: 'main'"
// @Param cron_id query int false "Type: '1'"
// @Param project_id query string false "ProjectID: '1'"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription [delete]
func (o *subscriptionController) DeleteAll(c *gin.Context) {
	var query = m.SubscribeQueryDto{Page: -1}
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
		return
	}

	model, err := o.service.Subscription.Read(&query)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	var items = 0
	for _, item := range model {
		if err := o.service.Cron.Delete(&m.CronQueryDto{ID: item.CronID}); err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		num, err := o.service.Subscription.Delete(&m.SubscribeQueryDto{CronID: item.CronID})
		if err != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
			return
		}

		items += num
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []string{},
		Items:  items,
	})
}

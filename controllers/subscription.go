package controllers

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/interfaces"
	"api/logs"
	m "api/models"
	"api/service"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type subscriptionController struct {
	service *service.SubscriptionService
}

func NewSubscriptionController(s *service.SubscriptionService) interfaces.Default {
	return &subscriptionController{service: s}
}

// @Tags Subscription
// @Summary Create list of Subscriptions to run operation
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "ID: '1'"
// @Param name query string false "Name: 'CodeRain'"
// @Param namespace query string false "Namespace: 'test'"
// @Param model body m.SubscribeDto true "Small info about subscription for k3s"
// @Success 201 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/list [post]
func (*subscriptionController) CreateAll(c *gin.Context) {

}

// @Tags Subscription
// @Summary Create Subscription to run operation
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "ID: '1'"
// @Param name query string false "Name: 'CodeRain'"
// @Param namespace query string false "Namespace: 'test'"
// @Param model body m.SubscribeDto true "Small info about subscription for k3s"
// @Success 201 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription [post]
func (o *subscriptionController) CreateOne(c *gin.Context) {
	var body m.SubscribeDto
	var id = helper.GetID(c)

	if err := c.ShouldBind(&body); err != nil || !body.IsOK() || id == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t, id: %t }", body.IsOK(), id != 0))
		return
	}

	var model = m.Subscription{Name: body.Name, CronTime: body.CronTime}
	if err := o.service.Create(&model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.ResHandler(c, http.StatusCreated, &m.Success{
		Status: "OK",
		Result: []m.Link{model},
		Items:  1,
	})

	// handler, ok := config.GetOperation(body.Operation)
	// if !ok {
	// 	helper.ErrHandler(c, http.StatusNotFound, fmt.Sprintf("Operation '%s' not founded", body.Operation))
	// 	return
	// }

	// var path string
	// if path, err = helper.FormPathFromHandler(c, handler); err != nil {
	// 	helper.ErrHandler(c, http.StatusNotFound, err.Error())
	// 	return
	// }

	// hasher := md5.New()
	// hasher.Write([]byte(strconv.Itoa(rand.Intn(1000000) + 5000)))
	// token := hex.EncodeToString(hasher.Sum(nil))

	// var reqBody []byte
	// if reqBody, err = json.Marshal(&m.CronCreateDto{
	// 	CronTime: body.CronTime,
	// 	URL:      config.ENV.URL + path,
	// 	Method:   handler.Method,
	// 	Token:    token,
	// }); err != nil {
	// 	fmt.Printf("Ohh noo; Anyway: %v", err)
	// 	return
	// }

	// hasher = md5.New()
	// salt := strconv.Itoa(rand.Intn(1000000) + 5000)
	// hasher.Write([]byte(salt + config.ENV.BotKey))

	// var req *http.Request
	// if req, err = http.NewRequest("POST", config.ENV.BotUrl+"/cron/subscribe?key="+hex.EncodeToString(hasher.Sum(nil)), bytes.NewBuffer(reqBody)); err != nil {
	// 	fmt.Printf("Ohh noo; Anyway: %v", err)
	// 	return
	// }

	// req.Header.Set("X-Custom-Header", salt)
	// req.Header.Set("Content-Type", "application/json")

	// var res *http.Response
	// client := &http.Client{}
	// res, err = client.Do(req)
	// if err != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong response")
	// 	return
	// }

	// defer res.Body.Close()

	// if res.StatusCode != http.StatusCreated {
	// 	helper.ErrHandler(c, res.StatusCode, "Bot request error")
	// 	return
	// }

	// var cron = m.CronEntity{}
	// if err = json.NewDecoder(res.Body).Decode(&cron); err != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong response")
	// 	return
	// }

	// model := m.Subscription{
	// 	CronID:   cron.ID,
	// 	CronTime: cron.Exec.CronTime,
	// 	Method:   handler.Method,
	// 	Path:     path,
	// }
	// result := db.DB.Create(&model)

	// if result.Error != nil || result.RowsAffected == 0 {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
	// 	go logs.DefaultLog("/controllers/subscription.go", result.Error)
	// 	return
	// }

	// go func() {
	// 	hasher = md5.New()
	// 	hasher.Write([]byte(token))

	// 	ctx := context.Background()
	// 	db.Redis.Set(ctx, "TOKEN:"+hex.EncodeToString(hasher.Sum(nil)), "OK", 0)
	// }()

	// go db.FlushValue("SUBSCRIPTION")
	// helper.ResHandler(c, http.StatusCreated, &m.Success{
	// 	Status: "OK",
	// 	Result: []m.Subscription{model},
	// 	Items:  1,

	// 	// TODO: Maybe on day I'll add this ....
	// 	// TotalItems: items,
	// })
}

func (*subscriptionController) ReadAll(c *gin.Context) {}

// @Tags Subscription
// @Summary Read subscription by id/cron_id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "This id can be a ID (Primary Key) or a CronID"
// @Success 200 {object} m.Success{result=[]m.Subscription}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/{id} [get]
func (*subscriptionController) ReadOne(c *gin.Context) {
	var id string
	var model []m.Subscription

	if id = c.Param("id"); id == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	var query = "cron_id = ?"
	if _, err := strconv.Atoi(id); err == nil {
		query = "id = ?"
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("ID=%s", id)))
	if err := helper.PrecacheResult(fmt.Sprintf("SUBSCRIPTION:%s", hex.EncodeToString(hasher.Sum(nil))), db.DB.Where(query, id), &model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/subscription.go", err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: model,
		Items:  1,

		// TODO: Maybe one day ....
		// TotalItems: items,
	})
}

func (*subscriptionController) UpdateOne(c *gin.Context) {}
func (*subscriptionController) UpdateAll(c *gin.Context) {}

// @Tags Subscription
// @Summary Delete subscription by id/cron_id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path string true "This id can be a ID (Primary Key) or a CronID"
// @Success 200 {object} m.Success{result=[]string{}}
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /subscription/{id} [delete]
func (*subscriptionController) DeleteOne(c *gin.Context) {
	var id string
	var model []m.Subscription

	if id = c.Param("id"); id == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	var query = "cron_id = ?"
	if _, err := strconv.Atoi(id); err == nil {
		query = "id = ?"
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("ID=%s", id)))
	if err := helper.PrecacheResult(fmt.Sprintf("SUBSCRIPTION:%s", hex.EncodeToString(hasher.Sum(nil))), db.DB.Where(query, id), &model); err != nil || len(model) == 0 {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/subscription.go", err.Error())
		return
	}

	hasher = md5.New()
	salt := strconv.Itoa(rand.Intn(1000000) + 5000)
	hasher.Write([]byte(salt + config.ENV.BotKey))

	var req *http.Request
	var err error
	if req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/cron/subscribe?key=%s&id=%s", config.ENV.BotUrl, hex.EncodeToString(hasher.Sum(nil)), model[0].CronID), nil); err != nil {
		fmt.Printf("Ohh noo; Anyway: %v", err)
		return
	}

	req.Header.Set("X-Custom-Header", salt)
	req.Header.Set("Content-Type", "application/json")

	var res *http.Response
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong response")
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		helper.ErrHandler(c, res.StatusCode, "Bot request error")
		return
	}

	db.DB.Where("id = ?", model[0].ID).Delete(&m.Subscription{})
	go db.FlushValue("SUBSCRIPTION")
	helper.ResHandler(c, http.StatusOK, &m.Success{
		Status: "OK",
		Result: []string{},
	})
}

func (*subscriptionController) DeleteAll(c *gin.Context) {}

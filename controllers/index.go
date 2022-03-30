package controllers

import (
	"api/helper"
	"api/interfaces"
	"api/logs"
	m "api/models"
	"api/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type indexController struct {
	service *service.IndexService
}

func NewIndexController(s *service.IndexService) interfaces.Index {
	return &indexController{service: s}
}

// @Summary Ping/Pong
// @Accept json
// @Produce application/json
// @Success 200 {object} m.Ping
// @failure 429 {object} m.Error
// @Router /ping [get]
func (*indexController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, m.Ping{
		Status:  "OK",
		Message: "pong",
	})
}

// @Summary Trace Ip :ip
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param ip path string true "Client IP"
// @Success 200 {object} m.Success{result=[]m.GeoIpLocations}
// @failure 429 {object} m.Error
// @failure 400 {object} m.Error
// @failure 500 {object} m.Error
// @Router /trace/{ip} [get]
func (o *indexController) TraceIp(c *gin.Context) {
	var ip string
	if ip = c.Param("ip"); ip == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect ip value")
		return
	}

	models, err := o.service.TraceIP(ip)
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

// @Summary Login
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model body m.LoginDto true "Login info"
// @Success 200 {object} m.TokenEntity
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /login [post]
func (o *indexController) Login(c *gin.Context) {
	var login m.LoginDto
	if err := c.ShouldBind(&login); err != nil || !login.IsOK() || len(strings.Split(login.Pass, "$")) != 2 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body")
		return
	}

	token, err := o.service.Login(&login)
	if err != nil {
		helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
		go logs.SendLogs(&m.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/refresh",
			File:    "/controllers/index.go",
			Message: "Something went wrong with auth",
			Desc:    err.Error(),
		})
		return
	}

	helper.ResHandler(c, http.StatusOK, m.TokenEntity{
		Status:       "OK",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

// @Summary Refresh access token
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Success 200 {object} m.TokenEntity
// @failure 400 {object} m.Error
// @failure 401 {object} m.Error
// @failure 422 {object} m.Error
// @failure 429 {object} m.Error
// @failure 500 {object} m.Error
// @Router /refresh [post]
func (o *indexController) Refresh(c *gin.Context) {
	var body m.TokenDto
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Refresh token not specified")
		return
	}

	token, err := o.service.Refresh(&body)
	if err != nil {
		helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
		go logs.SendLogs(&m.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/refresh",
			File:    "/controllers/index.go",
			Message: "Something went wrong with auth",
			Desc:    err.Error(),
		})
		return
	}

	helper.ResHandler(c, http.StatusOK, m.TokenEntity{
		Status:       "OK",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

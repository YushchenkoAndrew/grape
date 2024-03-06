package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthT interface {
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

type authController struct {
	service *authService
}

func NewAuthController(s *authService) AuthT {
	return &authController{service: s}
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
func (o *authController) Login(c *gin.Context) {
	// var login m.LoginDto
	// if err := c.ShouldBind(&login); err != nil || !login.IsOK() || len(strings.Split(login.Pass, "$")) != 2 {
	// 	helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body")
	// 	return
	// }

	// token, err := o.service.Login(&login)
	// if err != nil {
	// 	helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
	// 	go logs.SendLogs(&logs.Message{
	// 		Stat:    "ERR",
	// 		Name:    "grape",
	// 		Url:     "/api/refresh",
	// 		File:    "/controllers/index.go",
	// 		Message: "Something went wrong with auth",
	// 		Desc:    err.Error(),
	// 	})
	// 	return
	// }

	// helper.ResHandler(c, http.StatusOK, m.TokenEntity{
	// 	Status:       "OK",
	// 	AccessToken:  token.AccessToken,
	// 	RefreshToken: token.RefreshToken,
	// })
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
func (o *authController) Refresh(c *gin.Context) {
	// var body m.TokenDto
	// if err := c.ShouldBind(&body); err != nil {
	// 	helper.ErrHandler(c, http.StatusBadRequest, "Refresh token not specified")
	// 	return
	// }

	// token, err := o.service.Refresh(&body)
	// if err != nil {
	// 	helper.ErrHandler(c, http.StatusUnauthorized, err.Error())
	// 	go logs.SendLogs(&logs.Message{
	// 		Stat:    "ERR",
	// 		Name:    "grape",
	// 		Url:     "/api/refresh",
	// 		File:    "/controllers/index.go",
	// 		Message: "Something went wrong with auth",
	// 		Desc:    err.Error(),
	// 	})
	// 	return
	// }

	// helper.ResHandler(c, http.StatusOK, m.TokenEntity{
	// 	Status:       "OK",
	// 	AccessToken:  token.AccessToken,
	// 	RefreshToken: token.RefreshToken,
	// })
}

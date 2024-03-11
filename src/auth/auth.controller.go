package auth

import (
	"grape/src/auth/dto/request"
	"grape/src/common/dto/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *AuthService
}

func NewAuthController(s *AuthService) *AuthController {
	return &AuthController{service: s}
}

// @Summary Login
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model body request.LoginDto true "Login info"
// @Success 200 {object} response.LoginResponseDto
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var body request.LoginDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Login(&body)
	if err != nil {
		response.ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.Build(ctx, http.StatusOK, res)
}

// @Summary Refresh access token
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model body request.RefreshDto true "Login info"
// @Success 200 {object} response.LoginResponseDto
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /refresh [post]
func (c *AuthController) Refresh(ctx *gin.Context) {
	var body request.RefreshDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.ThrowErr(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Refresh(&body)
	if err != nil {
		response.ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.Build(ctx, http.StatusOK, res)
}

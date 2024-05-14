package auth

import (
	"grape/src/auth/dto/request"
	t "grape/src/auth/types"
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

// @Tags Auth
// @Summary Login
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model body request.LoginDto true "Login info"
// @Success 200 {object} response.LoginResponseDto
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var body request.LoginDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Login(&body)
	if err != nil {
		response.ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.Build(ctx, http.StatusOK, res)
}

// @Tags Auth
// @Summary Refresh access token
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param model body request.RefreshDto true "Login info"
// @Success 200 {object} response.LoginResponseDto
// @failure 400 {object} response.Error
// @failure 422 {object} response.Error
// @Router /auth/refresh [post]
func (c *AuthController) Refresh(ctx *gin.Context) {
	var body request.RefreshDto
	if err := ctx.ShouldBind(&body); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	res, err := c.service.Refresh(&body)
	if err != nil {
		response.ThrowErr(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}

	response.Build(ctx, http.StatusOK, res)
}

// @Tags Auth
// @Summary Logout
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Success 201 {object} interface{}
// @failure 401 {object} response.Error
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	claim, _ := ctx.Get("access_claim")
	res, err := c.service.Logout(claim.(*t.AccessClaim))

	response.Handler(ctx, http.StatusOK, res, err)
}

// @Tags Auth
// @Summary Check if jwt is valid
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Success 200 {object} interface{}
// @failure 401 {object} response.Error
// @Router /auth/ping [get]
func (c *AuthController) Ping(ctx *gin.Context) {
	response.Handler(ctx, http.StatusOK, gin.H{"message": "pong"}, nil)
}

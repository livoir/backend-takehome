package http

import (
	"app/domain"
	"app/pkg/common"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthUseCase domain.AuthUseCase
}

func NewAuthHandler(r *gin.RouterGroup, authUseCase domain.AuthUseCase) {
	handler := &AuthHandler{
		AuthUseCase: authUseCase,
	}
	r.POST("/login", handler.Login)
	r.POST("/register", handler.Register)

}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var request *domain.LoginRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	response, err := h.AuthUseCase.Login(ctx, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.SetCookie("AUTHORIZATION", response.AccessToken, 300, "/", "", false, true)
	ctx.SetCookie("REFRESH_TOKEN", response.RefreshToken, 3600*24*7, "/", "", false, true)
	handleOK(ctx, response)
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var request *domain.RegisterRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	response, err := h.AuthUseCase.Register(ctx, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOK(ctx, response)
}

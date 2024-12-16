package http

import (
	"app/domain"
	"app/pkg/common"
	"app/pkg/logger"

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
		logger.Log.Error(err.Error())
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	response, err := h.AuthUseCase.Login(ctx, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.SetCookie("AUTHORIZATION", response.AccessToken, 0, "/", "", false, true)
	ctx.SetCookie("REFRESH_TOKEN", response.RefreshToken, 0, "/", "", false, true)
	handleOK(ctx, response)
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var request *domain.RegisterRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Log.Error(err.Error())
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	response, err := h.AuthUseCase.Register(ctx, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOKCreated(ctx, response)
}

package http

import (
	"app/domain"
	"app/pkg/common"
	"app/pkg/logger"
	"net/http"
	"net/mail"
	"strings"

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
	r.POST("/refresh-token", handler.RefreshToken)

}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var request *domain.LoginRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Log.Error(err.Error())
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
		handleError(ctx, err)
		return
	}
	if err := isValidEmail(request.Email); err != nil {
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
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
		handleError(ctx, err)
		return
	}
	if err := isValidEmail(request.Email); err != nil {
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

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	tokenCookie, err := ctx.Request.Cookie("REFRESH_TOKEN")
	if err != nil {
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
		handleError(ctx, err)
		ctx.Abort()
	}
	token := tokenCookie.Value
	if token == "" {
		err := common.ErrInvalidToken
		handleError(ctx, err)
		ctx.Abort()
	}
	response, err := h.AuthUseCase.RefreshToken(ctx, token)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.SetCookie("AUTHORIZATION", response.AccessToken, 0, "/", "", false, true)
	handleOK(ctx, response)
}

func isValidEmail(email string) error {
	if strings.Contains(email, " ") {
		return common.NewCustomError(http.StatusBadRequest, "email address should not contain space")
	}
	if strings.Contains(email, "<") || strings.Contains(email, ">") {
		return common.NewCustomError(http.StatusBadRequest, "email address should not contain < or >")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return common.NewCustomError(http.StatusBadRequest, "invalid email address")
	}
	return nil
}

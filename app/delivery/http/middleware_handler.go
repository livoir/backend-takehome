package http

import (
	"app/domain"
	"app/pkg/common"

	"github.com/gin-gonic/gin"
)

type MiddlewareHandler struct {
	authUsecase domain.AuthUseCase
}

func NewMiddlewareHandler(authUsecase domain.AuthUseCase) *MiddlewareHandler {
	return &MiddlewareHandler{
		authUsecase: authUsecase,
	}
}

func (h *MiddlewareHandler) AuthMiddleware(ctx *gin.Context) {
	tokenCookie, err := ctx.Request.Cookie("AUTHORIZATION")
	if err != nil {
		err := common.ErrInvalidToken
		handleError(ctx, err)
		ctx.Abort()
	}
	token := tokenCookie.Value
	if token == "" {
		err := common.ErrInvalidToken
		handleError(ctx, err)
		ctx.Abort()
	}
	res, err := h.authUsecase.VerifyToken(ctx, token)
	if err != nil {
		handleError(ctx, err)
		ctx.Abort()
	}
	ctx.Set("userID", res.UserID)
	ctx.Next()
}

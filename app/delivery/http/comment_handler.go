package http

import (
	"app/domain"
	"app/pkg/common"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentUseCase domain.CommentUsecase
}

func NewCommentHandler(r *gin.RouterGroup, middleware *MiddlewareHandler, commentUseCase domain.CommentUsecase) {
	handler := &CommentHandler{
		commentUseCase: commentUseCase,
	}
	r.GET("", handler.FindCommentsByPostID)

	r.Use(middleware.AuthMiddleware)
	r.POST("", handler.CreateComment)
}

func (h *CommentHandler) CreateComment(ctx *gin.Context) {
	var request domain.CreateCommentRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	var path struct {
		PostID int64 `uri:"postID" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	userID := ctx.GetInt64("userID")
	request.AuthorID = userID
	response, err := h.commentUseCase.CreateComment(ctx, path.PostID, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOKCreated(ctx, response)
}

func (h *CommentHandler) FindCommentsByPostID(ctx *gin.Context) {
	var request domain.SearchParam
	if err := ctx.ShouldBindQuery(&request); err != nil {
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	var path struct {
		PostID int64 `uri:"postID" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	comments, total, err := h.commentUseCase.FindCommentsByPostID(ctx, path.PostID, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handlePagination(ctx, comments, request.Page, request.Limit, total)
}

package http

import (
	"app/domain"
	"app/pkg/common"
	"app/pkg/logger"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postUseCase domain.PostUseCase
}

func NewPostHandler(r *gin.RouterGroup, middleware *MiddlewareHandler, postUseCase domain.PostUseCase) {
	handler := &PostHandler{
		postUseCase: postUseCase,
	}
	r.GET("/:id", handler.GetByID)
	r.GET("", handler.GetAll)

	// Apply middleware
	r.Use(middleware.AuthMiddleware)

	r.POST("", handler.Create)
	r.PUT("/:id", handler.Update)
	r.DELETE("/:id", handler.Delete)
}

func (h *PostHandler) Create(ctx *gin.Context) {
	var request *domain.CreatePostRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Log.Error(err.Error())
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	id := ctx.GetInt64("userID")
	request.AuthorID = id
	response, err := h.postUseCase.Create(ctx, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOK(ctx, response)
}

func (h *PostHandler) Update(ctx *gin.Context) {
	var request *domain.UpdatePostRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Log.Error(err.Error())
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	var path struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		logger.Log.Error(err.Error())
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	request.AuthorID = ctx.GetInt64("userID")
	response, err := h.postUseCase.Update(ctx, path.ID, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOK(ctx, response)
}

func (h *PostHandler) Delete(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	request := &domain.DeletePostRequestDTO{
		AuthorID: userID,
	}
	var path struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		logger.Log.Error(err.Error())
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	err := h.postUseCase.Delete(ctx, path.ID, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOK(ctx, nil)
}

func (h *PostHandler) GetByID(ctx *gin.Context) {
	var path struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		logger.Log.Error(err.Error())
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	response, err := h.postUseCase.GetByID(ctx, path.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOK(ctx, response)
}

func (h *PostHandler) GetAll(ctx *gin.Context) {
	var search domain.SearchParam
	if err := ctx.ShouldBindQuery(&search); err != nil {
		logger.Log.Error(err.Error())
		err = common.ErrInvalidParam
		handleError(ctx, err)
		return
	}
	response, err := h.postUseCase.GetAll(ctx, search)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOK(ctx, response)
}

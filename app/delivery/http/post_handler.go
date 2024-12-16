package http

import (
	"app/domain"
	"app/pkg/common"
	"app/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postUseCase domain.PostUseCase
}

func NewPostHandler(r *gin.RouterGroup, middleware *MiddlewareHandler, postUseCase domain.PostUseCase) {
	handler := &PostHandler{
		postUseCase: postUseCase,
	}
	r.GET("/:postID", handler.GetByID)
	r.GET("", handler.GetAll)

	// Apply middleware
	r.Use(middleware.AuthMiddleware)

	r.POST("", handler.Create)
	r.PUT("/:postID", handler.Update)
	r.DELETE("/:postID", handler.Delete)
}

func (h *PostHandler) Create(ctx *gin.Context) {
	var request *domain.CreatePostRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Log.Error(err.Error())
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
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
	handleOKCreated(ctx, response)
}

func (h *PostHandler) Update(ctx *gin.Context) {
	var request *domain.UpdatePostRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Log.Error(err.Error())
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
		handleError(ctx, err)
		return
	}
	var path struct {
		PostID int64 `uri:"postID" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		logger.Log.Error(err.Error())
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
		handleError(ctx, err)
		return
	}
	request.AuthorID = ctx.GetInt64("userID")
	response, err := h.postUseCase.Update(ctx, path.PostID, request)
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
		PostID int64 `uri:"postID" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		logger.Log.Error(err.Error())
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
		handleError(ctx, err)
		return
	}
	err := h.postUseCase.Delete(ctx, path.PostID, request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleOK(ctx, nil)
}

func (h *PostHandler) GetByID(ctx *gin.Context) {
	var path struct {
		PostID int64 `uri:"postID" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		logger.Log.Error(err.Error())
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
		handleError(ctx, err)
		return
	}
	response, err := h.postUseCase.GetByID(ctx, path.PostID)
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
		err = common.NewCustomError(http.StatusBadRequest, err.Error())
		handleError(ctx, err)
		return
	}
	if search.Limit == 0 {
		search.Limit = 10
	}
	if search.Page == 0 {
		search.Page = 1
	}
	response, total, err := h.postUseCase.GetAll(ctx, search)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handlePagination(ctx, response, search.Page, search.Limit, total)
}

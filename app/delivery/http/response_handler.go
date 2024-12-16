package http

import (
	"app/pkg/common"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Success   bool        `json:"success"`
	Status    int         `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type PagedResponse struct {
	Success   bool        `json:"success"`
	Status    int         `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Page      int         `json:"page"`
	Size      int         `json:"size"`
	Total     int64       `json:"total"`
	Data      interface{} `json:"data"`
}

func handleOKCreated(ctx *gin.Context, data interface{}) {
	now := time.Now()
	ctx.JSON(http.StatusCreated, BaseResponse{
		Success:   true,
		Status:    http.StatusCreated,
		Timestamp: now,
		Data:      data,
	})
}

func handleOK(ctx *gin.Context, data interface{}) {
	now := time.Now()
	ctx.JSON(http.StatusOK, BaseResponse{
		Success:   true,
		Status:    http.StatusOK,
		Timestamp: now,
		Data:      data,
	})
}

func handlePagination(ctx *gin.Context, data interface{}, page, size int, total int64) {
	now := time.Now()
	ctx.JSON(http.StatusOK, PagedResponse{
		Success:   true,
		Status:    http.StatusOK,
		Timestamp: now,
		Page:      page,
		Size:      size,
		Total:     total,
		Data:      data,
	})
}

func handleError(ctx *gin.Context, err error) {
	now := time.Now()
	if customErr, ok := err.(*common.CustomError); ok {
		errData := gin.H{"success": false, "timestamp": now, "message": customErr.Message}
		ctx.JSON(customErr.StatusCode, errData)
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "timestamp": now, "message": "internal server error"})
}

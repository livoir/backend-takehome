package http

import (
	"app/pkg/common"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handleOK(ctx *gin.Context, data interface{}) {
	now := time.Now()
	ctx.JSON(http.StatusOK, gin.H{"success": true, "timestamp": now, "data": data})
}

func handlePagination(ctx *gin.Context, data interface{}, page, size, total int) {
	now := time.Now()
	ctx.JSON(http.StatusOK, gin.H{"success": true, "timestamp": now, "data": data, "page": page, "size": size, "total": total})
}

func handleError(ctx *gin.Context, err error) {
	now := time.Now()
	if customErr, ok := err.(*common.CustomError); ok {
		errData := gin.H{"success": false, "timestamp": now, "message": customErr.Message}
		switch customErr.StatusCode {
		case http.StatusNotFound:
			ctx.JSON(http.StatusNotFound, errData)
		case http.StatusForbidden:
			ctx.JSON(http.StatusForbidden, errData)
		case http.StatusConflict:
			ctx.JSON(http.StatusConflict, errData)
		case http.StatusBadRequest:
			ctx.JSON(http.StatusBadRequest, errData)
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "timestamp": now, "message": "internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "timestamp": now, "message": "internal server error"})
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (e *EventLayer) CheckUsers(ctx *gin.Context) {
	start, ok := ctx.GetQuery("start")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "can't get query start"})
		ctx.Abort()
		return
	}

	end, ok := ctx.GetQuery("stop")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "can't get query stop"})
		ctx.Abort()
		return
	}
	intStart, err := strconv.ParseInt(start, 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "error while asserting 'start'"})
		ctx.Abort()
		return
	}

	intEnd, err := strconv.ParseInt(end, 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "error while asserting 'end'"})
		ctx.Abort()
		return
	}

	users, err := e.App.CheckUsers(int(intStart), int(intEnd))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "err": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "users": users})
}

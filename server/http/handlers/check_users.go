package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (e *EventLayer) CheckUsers(ctx *gin.Context) {
	var intStart int64
	var intEnd int64
	var err error
	if start, ok := ctx.GetQuery("start"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "can't get query start"})
		ctx.Abort()
		return
	} else {
		intStart, err = strconv.ParseInt(start, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "error while asserting 'start'"})
			ctx.Abort()
			return
		}
	}

	if end, ok := ctx.GetQuery("stop"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "can't get query stop"})
		ctx.Abort()
		return
	} else {
		intEnd, err = strconv.ParseInt(end, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "error while asserting 'end'"})
			ctx.Abort()
			return
		}
	}

	users, err := e.App.CheckUsers(intStart, intEnd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "err": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "users": users})
}

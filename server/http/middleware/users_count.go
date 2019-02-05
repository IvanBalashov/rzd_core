package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *EventLayer) UsersCount(ctx *gin.Context) {
	users, err := e.App.GetUsersList()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "length": len(users)})
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *EventLayer) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

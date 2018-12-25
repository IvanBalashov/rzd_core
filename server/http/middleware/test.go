package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *AppLayer) Test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}


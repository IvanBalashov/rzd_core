package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *AppLayer) Test(ctx *gin.Context) {
	a.App.GetSeats([]int{0,1})
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

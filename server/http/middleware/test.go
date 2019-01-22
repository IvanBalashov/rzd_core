package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/entity"
)

func (a *AppLayer) Test(ctx *gin.Context) {
	a.App.GetSeats(entity.RouteArgs{
		Dir:          "0",
		Tfl:          "3",
		Code0:        "2000000",
		Code1:        "2010000",
		Dt0:          "25.01.2019",
		CheckSeats:   "0",
		WithOutSeats: "y",
		Version:      "v.2018",
	})
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

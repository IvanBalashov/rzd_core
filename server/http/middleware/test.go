package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/entity"
	"strconv"
)

func (a *AppLayer) Test(ctx *gin.Context) {
	code1, code2, err := a.App.GetCodes("Москва", "ЯРОСЛАВЛЬ-ГЛАВНЫЙ")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	err = a.App.GetSeats(entity.RouteArgs{
		Dir:          "0",
		Tfl:          "3",
		Code0:        strconv.Itoa(code1),
		Code1:        strconv.Itoa(code2),
		Dt0:          "25.01.2019",
		CheckSeats:   "0",
		WithOutSeats: "y",
		Version:      "v.2018",
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

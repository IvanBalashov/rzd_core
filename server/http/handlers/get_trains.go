package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/entity"
	"rzd/server"
	"strconv"
)

type SeatsArgs struct {
	Direction string `form:"dir" binding:"required"`
	Target    string `form:"target" binding:"required"`
	Source    string `form:"source" binding:"required"`
	Date      string `form:"date" binding:"required"`
}

func (e *EventLayer) GetAllTrains(ctx *gin.Context) {
	query := &SeatsArgs{}

	if err := ctx.ShouldBindQuery(query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		ctx.Abort()
		return
	}

	code1, code2, err := e.App.GetStationCodes(query.Target, query.Source)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		ctx.Abort()
		return
	}

	getTrainsReq := &entity.RouteArgs{
		Dir:          query.Direction,
		Tfl:          "1",
		Code0:        strconv.Itoa(code1),
		Code1:        strconv.Itoa(code2),
		Dt0:          query.Date,
		CheckSeats:   "0",
		WithOutSeats: "y",
		Version:      "v.2018",
	}
	routes, err := e.App.GetInfoAboutTrains(getTrainsReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": err.Error()})
		ctx.Abort()
		return
	}

	trains := []server.Trains{}
	seats := map[string]server.Seats{}
	for _, val := range routes {
		for k, v := range val.Seats {
			seats[string(k)] = server.Seats{
				Count: v.SeatsCount,
				Price: v.Price,
				Chosen: v.Chosen,
			}
		}
		trains = append(trains, server.Trains{
			TrainID:   val.ID,
			MainRoute: val.Route0 + " - " + val.Route1,
			Segment:   val.Station + " - " + val.Station1,
			StartDate: val.Date0 + "_" + val.Time0,
			EndTime:   val.Date1 + "_" + val.Time1,
			Seats:     seats,
		})
		seats = map[string]server.Seats{}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "trains": trains})
}

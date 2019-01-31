package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/entity"
	"strconv"
)

type SeatsArgs struct {
	Direction string `form:"dir" binding:"required"`
	Target    string `form:"target" binding:"required"`
	Source    string `form:"source" binding:"required"`
	Date      string `form:"date" binding:"required"`
}

func (a *AppLayer) GetAllTrains(ctx *gin.Context) {
	query := SeatsArgs{}
	trains := []Trains{}
	seats := []Seats{}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	code1, code2, err := a.App.GetStationCodes(query.Target, query.Source)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}

	routes, err := a.App.GetInfoAboutTrains(entity.RouteArgs{
		Dir:          query.Direction,
		Tfl:          "1",
		Code0:        strconv.Itoa(code1),
		Code1:        strconv.Itoa(code2),
		Dt0:          query.Date,
		CheckSeats:   "0",
		WithOutSeats: "y",
		Version:      "v.2018",
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}

	// Parsing answer here coz we need one answer for all "servres"
	for _, val := range routes {
		for i := range val.Seats {
			seats = append(seats, Seats{
				Name:  val.Seats[i].SeatsName,
				Count: val.Seats[i].SeatsCount,
				Price: val.Seats[i].Price,
			})
		}
		trains = append(trains, Trains{
			TrainID:   val.ID,
			MainRoute: val.Route0 + " - " + val.Route1,
			Segment:   val.Station + " - " + val.Station1,
			StartDate: val.Date0 + "_" + val.Time0,
			EndTime:   val.Date1 + "_" + val.Time1,
			Seats:     seats,
		})
		seats = []Seats{}
	}
	ctx.JSON(http.StatusOK, trains)
}

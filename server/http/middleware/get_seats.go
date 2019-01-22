package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rzd/app/entity"
	"strconv"
)

//TODO: parse query in to args.
func (a *AppLayer) GetSeats(ctx *gin.Context) {
	trains := []Trains{}
	seats := []Seats{}
	code1, code2, err := a.App.GetCodes("Москва", "ЯРОСЛАВЛЬ")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	routes, err := a.App.GetSeats(entity.RouteArgs{
		Dir:          "0",
		Tfl:          "1",
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
			MainRoute: val.Route0 + " - " + val.Route0,
			Segment:   val.Station + " - " + val.Station1,
			StartDate: val.Date0 + "_" + val.Time0,
			Seats:     seats,
		})
	}
	ctx.JSON(http.StatusOK, trains)
}

type Trains struct {
	MainRoute string  `json:"main_route"`
	Segment   string  `json:"segment"`
	StartDate string  `json:"start_date"`
	Seats     []Seats `json:"seats"`
}

type Seats struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Price int    `json:"price"`
}

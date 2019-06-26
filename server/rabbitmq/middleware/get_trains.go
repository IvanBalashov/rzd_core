package middleware

import (
	"encoding/json"
	"fmt"
	"rzd/app/entity"
	"strconv"
)

func (e *EventLayer) GetAllTrains(query interface{}) (interface{}, error) {
	request := AllTrainsRequest{}
	response := []Trains{}

	if data, err := json.Marshal(query); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(data, &request)
		if err != nil {
			return nil, err
		}
	}

	code1, code2, err := e.App.GetStationCodes(request.Target, request.Source)
	if err != nil {
		e.LogChanel <- fmt.Sprintf("RabbitMQ->GetInfoAboutTrains: Error in GetStationCodes - %s", err)
		return nil, err
	}

	routes, err := e.App.GetInfoAboutTrains(entity.RouteArgs{
		Dir:          request.Direction,
		Tfl:          "1",
		Code0:        strconv.Itoa(code1),
		Code1:        strconv.Itoa(code2),
		Dt0:          request.Date,
		CheckSeats:   "0",
		WithOutSeats: "y",
		Version:      "v.2018",
	})

	for _, val := range routes {
		seats := []Seats{}
		for i := range val.Seats {
			seats = append(seats, Seats{
				Name:  val.Seats[i].SeatsName,
				Count: val.Seats[i].SeatsCount,
				Price: val.Seats[i].Price,
			})
		}
		response = append(response, Trains{
			TrainID:   val.ID,
			MainRoute: fmt.Sprintf("%s-%s", val.Route0, val.Route0),
			Segment:   fmt.Sprintf("%s-%s", val.Station, val.Station1),
			StartDate: fmt.Sprintf("%s_%s", val.Date0, val.Time0),
			EndTime:   fmt.Sprintf("%s_%s", val.Date1, val.Time1),
			Seats:     seats,
		})
	}

	return response, nil
}

package middleware

import (
	"encoding/json"
	"fmt"
	"rzd/app/entity"
	"strconv"
)

func (m *EventLayer) GetAllTrains(query interface{}) (interface{}, error) {
	response := []Trains{}
	request := AllTrainsRequest{}

	if data, err := json.Marshal(query); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(data, &request)
		if err != nil {
			return nil, err
		}
	}

	code1, code2, err := m.App.GetStationCodes(request.Target, request.Source)
	if err != nil {
		m.LogChanel <- fmt.Sprintf("RabbitMQ->GetInfoAboutTrains: Error in GetStationCodes - %s", err)
		return nil, err
	}

	routes, err := m.App.GetInfoAboutTrains(entity.RouteArgs{
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
			MainRoute: val.Route0 + " - " + val.Route0,
			Segment:   val.Station + " - " + val.Station1,
			StartDate: val.Date0 + "_" + val.Time0,
			Seats:     seats,
		})
	}

	return response, nil
}

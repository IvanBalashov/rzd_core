package middleware

import (
	"encoding/json"
	"fmt"
	"rzd/app/entity"
	"rzd/server"
	"strconv"
)

func (e *EventLayer) GetAllTrains(query interface{}) (interface{}, error) {
	request := AllTrainsRequest{}
	response := []server.Trains{}

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
		seats := map[string]server.Seats{}
		for k, v := range val.Seats {
			seats[string(k)] = server.Seats{
				Count: v.SeatsCount,
				Price: v.Price,
				Chosen: v.Chosen,
			}
		}
		response = append(response, server.Trains{
			TrainID:   val.ID,
			MainRoute: fmt.Sprintf("%s-%s", val.Route0, val.Route1),
			Segment:   fmt.Sprintf("%s-%s", val.Station, val.Station1),
			StartDate: fmt.Sprintf("%s_%s", val.Date0, val.Time0),
			EndTime:   fmt.Sprintf("%s_%s", val.Date1, val.Time1),
			Seats:     seats,
		})
	}

	return response, nil
}

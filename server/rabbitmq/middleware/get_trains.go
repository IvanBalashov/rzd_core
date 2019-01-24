package middleware

import (
	"encoding/json"
	"fmt"
	"rzd/app/entity"
	"strconv"
)

func (m *EventLayer) GetSeats(query Data, eventName string) ([]byte, error) {
	event := GetAllTrainsEvent{}
	event.Event = eventName
	seats := []Seats{}
	code1, code2, err := m.App.GetCodes(query.Target, query.Source)
	if err != nil {
		m.LogChanel <- fmt.Sprintf("RabbitMQ->GetSeats: Error in GetCodes - %s", err)
		return nil, err
	}

	routes, err := m.App.GetSeats(entity.RouteArgs{
		Dir:          query.Direction,
		Tfl:          "1",
		Code0:        strconv.Itoa(code1),
		Code1:        strconv.Itoa(code2),
		Dt0:          query.Date,
		CheckSeats:   "0",
		WithOutSeats: "y",
		Version:      "v.2018",
	})

	// full logic like in http middleware, can be rewrited
	for _, val := range routes {
		for i := range val.Seats {
			seats = append(seats, Seats{
				Name:  val.Seats[i].SeatsName,
				Count: val.Seats[i].SeatsCount,
				Price: val.Seats[i].Price,
			})
		}
		event.Data = append(event.Data, Trains{
			MainRoute: val.Route0 + " - " + val.Route0,
			Segment:   val.Station + " - " + val.Station1,
			StartDate: val.Date0 + "_" + val.Time0,
			Seats:     seats,
		})
		seats = []Seats{}
	}
	data, err := json.Marshal(event)
	if err != nil {
		m.LogChanel <- fmt.Sprintf("RabbitMQ->GetSeats: Error in GetCodes - %s", err)
		return nil, err
	}
	return data, nil
}

package middleware

import (
	"encoding/json"
	"fmt"
	"rzd/app/entity"
	"strconv"
)

func (m *EventLayer) SaveInfoAboutTrain(query interface{}) (interface{}, error) {
	response := []Trains{}
	request := SaveOneTrainRequest{}
	seats := []Seats{}

	if data, err := json.Marshal(query); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(data, &request)
		if err != nil {
			return nil, err
		}
	}

	routes, err := m.App.SaveInfoAboutTrain()

	for _, val := range routes {
		for i := range val.Seats {
			seats = append(seats, Seats{
				Name:  val.Seats[i].SeatsName,
				Count: val.Seats[i].SeatsCount,
				Price: val.Seats[i].Price,
			})
		}
		response = append(response, Trains{
			MainRoute: val.Route0 + " - " + val.Route0,
			Segment:   val.Station + " - " + val.Station1,
			StartDate: val.Date0 + "_" + val.Time0,
			Seats:     seats,
		})
		seats = []Seats{}
	}

	return response, nil
}

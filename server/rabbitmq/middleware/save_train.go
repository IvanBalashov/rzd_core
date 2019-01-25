package middleware

import (
	"encoding/json"
)

func (m *EventLayer) SaveInfoAboutTrain(query interface{}) (interface{}, error) {
	//response := []Trains{}
	request := SaveOneTrainRequest{}
	//seats := []Seats{}

	if data, err := json.Marshal(query); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(data, &request)
		if err != nil {
			return nil, err
		}
	}

	err := m.App.SaveInfoAboutTrain()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

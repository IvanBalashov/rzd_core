package middleware

import (
	"encoding/json"
)

func (m *EventLayer) SaveInfoAboutTrain(query interface{}) (interface{}, error) {
	request := Trains{}

	if data, err := json.Marshal(query); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(data, &request)
		if err != nil {
			return nil, err
		}
	}

	err := m.App.SaveInfoAboutTrain(request.TrainID)
	if err != nil {
		return nil, err
	}

	return Status{Status: "OK"}, nil
}

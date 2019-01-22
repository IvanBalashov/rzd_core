package middleware

import (
	"encoding/json"
	"rzd/app/entity"
)

func (m *EventLayer) GetSeats(ids string) ([]byte, error) {
	trains := m.App.GetSeats(entity.RouteArgs{})
	data, err := json.Marshal(trains)
	if err != nil {
		return nil, err
	}
	return data, nil
}

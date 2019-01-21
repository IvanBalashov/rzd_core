package middleware

import "encoding/json"

func (m *EventLayer) GetSeats(ids string) ([]byte, error) {
	trains, err := m.App.GetSeats([]int{0, 1})
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(trains)
	if err != nil {
		return nil, err
	}
	return data, nil
}

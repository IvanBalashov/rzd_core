package middleware

func (m *EventLayer) SaveTrain(query interface{}) (interface{}, error) {
	m.App.SaveTrain()
	return nil, nil
}

package middleware

import "encoding/json"

func (e *EventLayer) CheckUsers(query interface{}) (interface{}, error) {
	request := CheckUsers{}

	if data, err := json.Marshal(query); err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(data, &request)
		if err != nil {
			return nil, err
		}
	}
	_, err := e.App.CheckUsers(request.Start, request.End)
	if err != nil {
		return Status{"error"}, err
	}
	return Status{"ok"}, nil
}

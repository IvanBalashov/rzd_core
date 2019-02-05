package middleware

import "encoding/json"

func (e *EventLayer) CheckUsers(query interface{}) (interface{}, error) {
	request := CheckUsersRequest{}

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
		return StatusResponse{"error"}, err
	}
	return StatusResponse{"ok"}, nil
}

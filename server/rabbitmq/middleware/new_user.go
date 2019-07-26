package middleware

import (
	"encoding/json"
	"rzd/app/entity"
)

func (e *EventLayer) NewUser(query interface{}) (interface{}, error) {
	request := &NewUserRequest{}
	if data, err := json.Marshal(query); err != nil {
		return StatusResponse{Status: "fail"}, err
	} else {
		err = json.Unmarshal(data, &request)
		if err != nil {
			return StatusResponse{Status: "fail"}, err
		}
	}

	user := &entity.User{
		UserTelegramID: request.UserTelegramID,
		UserName:       request.UserName,
		Notify:         request.Notify,
	}
	ok, err := e.App.AddUser(user)
	if err != nil {
		if ok {
			return StatusResponse{Status: "success"}, nil
		}
		return StatusResponse{Status: "fail"}, err
	}
	return StatusResponse{Status: "success"}, nil
}

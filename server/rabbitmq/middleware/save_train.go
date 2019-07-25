package middleware

import (
	"encoding/json"
	"rzd/app/entity"
)

func (e *EventLayer) SaveInfoAboutTrain(query interface{}) (interface{}, error) {
	request := &SaveOneTrainRequest{}

	if data, err := json.Marshal(query); err != nil {
		return StatusResponse{Status: "fail"}, err
	} else {
		err = json.Unmarshal(data, &request)
		if err != nil {
			return StatusResponse{Status: "fail"}, err
		}
	}

	trainID, err := e.App.SaveInfoAboutTrain(request.Train.TrainID)
	if err != nil {
		return StatusResponse{Status: "fail"}, err
	}

	user := &entity.User{
		UserTelegramID: request.User.UserID,
		UserName:       request.User.UserName,
	}
	err = e.App.SaveTrainInUser(user, trainID)
	if err != nil {
		return StatusResponse{Status: "fail"}, err
	}

	return StatusResponse{Status: "ok"}, nil
}

package middleware

import (
	"encoding/json"
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

	err = e.App.SaveTrainInUser(request.User.UserID, trainID)
	if err != nil {
		return StatusResponse{Status: "fail"}, err
	}

	return StatusResponse{Status: "ok"}, nil
}

package usecase

import (
	"encoding/json"
	"rzd/app/entity"
)

func (a *App) SaveInfoAboutTrain(trainID string) error {
	train := entity.Train{}

	data, err := a.Cache.Get(trainID)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &train)
	if err != nil {
		return err
	}

	err = a.Trains.Create(train)
	if err != nil {
		return err
	}

	return nil
}

package usecase

import (
	"errors"
	"fmt"
	"rzd/app/entity"
	"strconv"
)

func (a *App) GenerateTrainsList(route entity.Route) ([]entity.Train, error) {
	trains := []entity.Train{}
	newTrain := entity.Train{}

	if len(route.Tp) == 0 {
		return nil, errors.New(fmt.Sprintf("App->GenerateTrainsList: Got empty route array"))
	}

	for _, val := range route.Tp[0].List {
		seats := []entity.Seats{}
		for _, seatsInfo := range val.ServiceCategories {
			seats = append(seats, entity.Seats{
				SeatsCount: seatsInfo.FreeSeats,
				Price:      seatsInfo.Price,
				SeatsName:  seatsInfo.TypeLoc,
			})
		}
		newTrain = entity.Train{
			Number:   val.Number,
			Type:     strconv.Itoa(val.Type),
			Brand:    val.Brand,
			Route0:   val.Route0,
			Route1:   val.Route1,
			TrDate0:  val.TrDate0,
			TrTime0:  val.TrTime0,
			Station:  val.Station,
			Station1: val.Station1,
			Date0:    val.Date0,
			Time0:    val.Time0,
			Date1:    val.Date1,
			Time1:    val.Time1,
			Seats:    seats,
		}
		trains = append(trains, newTrain)
	}

	return trains, nil
}

package usecase

import (
	"rzd/app/entity"
	"strconv"
	"time"
)

// im think what need move here request for get rid and codes for trains.
func (a *App) GetInfoAboutTrains(args entity.RouteArgs) ([]entity.Train, error) {
	ridArgs := entity.RidArgs{
		Dir:          args.Dir,
		Tfl:          args.Tfl,
		CheckSeats:   args.CheckSeats,
		Code0:        args.Code0,
		Code1:        args.Code1,
		Dt0:          args.Dt0,
		WithOutSeats: args.WithOutSeats,
		Version:      args.Version,
	}

	rid, err := a.Routes.GetRid(ridArgs)
	if err != nil {
		a.LogChan <- err.Error()
		return nil, err
	}
	time.Sleep(750 * time.Millisecond)

	args.Rid = strconv.FormatInt(rid.RID, 10)
	route, err := a.Routes.GetRoutes(args)
	if err != nil {
		a.LogChan <- err.Error()
		return nil, err
	}

	trains, err := a.GenerateTrainsList(route)
	if err != nil {
		a.LogChan <- err.Error()
		return nil, err
	}

	return trains, nil
}

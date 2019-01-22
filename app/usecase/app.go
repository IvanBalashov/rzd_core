package usecase

import (
	"rzd/app/entity"
	"rzd/app/gateways/rzd_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"strconv"
	"time"
)

// TODO: Think about how correct work with error messages.
type App struct {
	Trains trains_gateway.TrainsGateway
	Users  users_gateway.UsersGateway
	Routes rzd_gateway.RzdGateway
}

func NewApp(trains trains_gateway.TrainsGateway, users users_gateway.UsersGateway, routes rzd_gateway.RzdGateway) App {
	return App{
		Trains: trains,
		Users:  users,
		Routes: routes,
	}
}

// im think what need move here request for get rid and codes for trains.
func (a *App) GetSeats(args entity.RouteArgs) ([]entity.Train, error) {
	ridArgs := entity.RidArgs{
		Dir:          args.Dir,
		Tfl:          args.Tfl,
		CheckSeats:   args.CheckSeats,
		Code0:        args.Code0,
		Code1:        args.Code1,
		Dt0:          args.Dt0,
		WithOutSeats: args.WithOutSeats,
		Version:      "v.2018", // FIXME: Now hardcoded, in future move this param in envs.
	}

	// cache for rid.
	rid, err := a.Routes.GetRid(ridArgs)
	if err != nil {
		return nil, err
	}
	args.Rid = strconv.FormatInt(rid.RID, 10)

	time.Sleep(time.Second)

	route, err := a.Routes.GetRoutes(args)
	if err != nil {
		return nil, err
	}

	trains, err := a.saveTrains(route)
	if err != nil {
		return nil, err
	}

	return trains, nil
}

func (a *App) saveTrains(route entity.Route) ([]entity.Train, error) {
	trains := []entity.Train{}
	newTrain := entity.Train{}
	for _, val := range route.Tp[0].List {
		seats := []entity.Seats{}
		for _, j := range val.ServiceCategories {
			seats = append(seats, entity.Seats{
				SeatsCount: j.FreeSeats,
				Price:      j.Price,
				SeatsName:  j.TypeLoc,
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

		err := a.Trains.Create(newTrain)
		trains = append(trains, newTrain)
		if err != nil {
			return nil, err
		}
	}
	return trains, nil
}

// TODO: can be parallel.
func (a *App) GetCodes(target, source string) (int, int, error) {
	var code1, code2 int
	code1, err := a.Routes.GetDirectionsCode(target)
	if err != nil {
		return 0, 0, err
	}

	code2, err = a.Routes.GetDirectionsCode(source)
	if err != nil {
		return 0, 0, err
	}

	return code1, code2, nil
}

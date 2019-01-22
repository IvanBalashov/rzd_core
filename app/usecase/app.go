package usecase

import (
	"fmt"
	"rzd/app/entity"
	"rzd/app/gateways/route_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"strconv"
	"time"
)

// TODO: Think about how correct work with error messages.
type App struct {
	Trains trains_gateway.TrainsGateway
	Users  users_gateway.UsersGateway
	Routes route_gateway.RzdGateway
}

func NewApp(trains trains_gateway.TrainsGateway, users users_gateway.UsersGateway, routes route_gateway.RzdGateway) App {
	return App{
		Trains: trains,
		Users:  users,
		Routes: routes,
	}
}

// im think what need move here request for get rid and codes for trains.
func (a *App) GetSeats(args entity.RouteArgs) error {
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

	rid, err := a.Routes.GetRid(ridArgs)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	args.Rid = strconv.FormatInt(rid.RID, 10)
	route, err := a.Routes.GetRoutes(args)
	if err != nil {
		return err
	}

	err = a.saveTrains(route)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) saveTrains(route entity.Route) error {
	fmt.Printf("route.TP[0].List %s\n", route.Tp[0].List)
	for _, val := range route.Tp[0].List {
		for _, j := range val.ServiceCategories {
			err := a.Trains.Create(entity.Train{
				Number:     val.Number,
				Type:       strconv.Itoa(val.Type),
				Brand:      val.Brand,
				Route0:     val.Route0,
				Route1:     val.Route1,
				TrDate0:    val.TrDate0,
				TrTime0:    val.TrTime0,
				Station:    val.Station,
				Station1:   val.Station1,
				Date0:      val.Date0,
				Time0:      val.Time0,
				Class:      j.TypeLoc,
				SeatsCount: strconv.FormatInt(int64(j.FreeSeats), 10),
				Price:      strconv.FormatInt(int64(j.Price), 10),
			})
			fmt.Printf("val - %+v\n", val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: can be parallel.
func (a *App) GetCodes(target, source string) (int, int, error) {
	code1, err := a.Routes.GetDirectionsCode(target)
	if err != nil {
		return 0, 0, err
	}
	code2, err := a.Routes.GetDirectionsCode(source)
	if err != nil {
		return 0, 0, err
	}
	return code1, code2, nil
}

package usecase

import (
	"fmt"
	"rzd/app/entity"
	"rzd/app/gateways/route_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"strconv"
)

// TODO: Think about how correct work with error messages.
type App struct {
	Trains trains_gateway.TrainsGateway
	Users  users_gateway.UsersGateway
	Routes route_gateway.RouteGateway
}

func NewApp(trains trains_gateway.TrainsGateway, users users_gateway.UsersGateway, routes route_gateway.RouteGateway) App {
	return App{
		Trains: trains,
		Users:  users,
		Routes: routes,
	}
}

func (a *App) GetSeats(args entity.RouteArgs) error {
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
				Type:       "0",
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
			fmt.Printf("val - %s\n", val)
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

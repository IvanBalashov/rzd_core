package usecase

import (
	"fmt"
	"rzd/app/entity"
	"rzd/app/gateways/route_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"strconv"
)

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

/*func (a *App) GetSeats(args entity.RouteArgs) error {
	route := a.Routes.GetRoutes(args)
	fmt.Printf("route - %s", route)
	err := a.saveTrains(route)
	if err != nil {
		return err
	}
	return nil
}*/
func (a *App) GetSeats(ids []int) ([]entity.Train, error) {
	return nil, nil
}

func (a *App) saveTrains(route entity.Route) error {
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

package usecase

import (
	"errors"
	"fmt"
	"rzd/app/entity"
	"rzd/app/gateways/rzd_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"strconv"
	"time"
)

// TODO: Think about how correct work with error messages.
type App struct {
	Trains  trains_gateway.TrainsGateway
	Users   users_gateway.UsersGateway
	Routes  rzd_gateway.RzdGateway
	LogChan chan string
	Rid     string
}

func NewApp(trains trains_gateway.TrainsGateway, users users_gateway.UsersGateway, routes rzd_gateway.RzdGateway, logChan chan string) App {
	return App{
		Trains:  trains,
		Users:   users,
		Routes:  routes,
		LogChan: logChan,
		Rid:     "",
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
		Version:      args.Version,
	}
	rid, err := a.Routes.GetRid(ridArgs)
	if err != nil {
		a.LogChan <- err.Error()
		return nil, err
	}
	a.Rid = strconv.FormatInt(rid.RID, 10)

	time.Sleep(500 * time.Millisecond)
	args.Rid = a.Rid
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

func (a *App) GetCodes(target, source string) (int, int, error) {
	var code1 = make(chan GoroutineAnswer)
	var answers = map[string]int{}
	go func() {
		data, err := a.Routes.GetDirectionsCode(target)
		if err != nil {
			a.LogChan <- err.Error()
		}
		code1 <- GoroutineAnswer{
			Code: data,
			Station: "target",
		}
	}()
	go func() {
		data, err := a.Routes.GetDirectionsCode(source)
		if err != nil {
			a.LogChan <- err.Error()
		}
		code1 <- GoroutineAnswer{
			Code: data,
			Station: "source",
		}
	}()
	for{
		select {
		case val, _ := <- code1:
			if val.Station == "target" {
				answers["target"] = val.Code
			} else {
				answers["source"] = val.Code
			}
		}
		if len(answers) == 2 {
			break
		}
	}
	return answers["target"], answers["source"], nil
}

type GoroutineAnswer struct {
	Code int
	Station string
}
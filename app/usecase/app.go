package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"rzd/app/entity"
	"rzd/app/gateways/cache_gateway"
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
	Cache   cache_gateway.CacheGateway
	LogChan chan string
}

func NewApp(trains trains_gateway.TrainsGateway, users users_gateway.UsersGateway, routes rzd_gateway.RzdGateway, cache cache_gateway.CacheGateway, logChan chan string) App {
	return App{
		Trains:  trains,
		Users:   users,
		Routes:  routes,
		Cache:   cache,
		LogChan: logChan,
	}
}

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

	rid, cookies, err := a.Routes.GetRid(ridArgs)
	if err != nil {
		a.LogChan <- err.Error()
		return nil, err
	}

	time.Sleep(450 * time.Millisecond)

	args.Rid = strconv.FormatInt(rid.RID, 10)
	route, err := a.Routes.GetRoutes(args, cookies)
	if err != nil {
		a.LogChan <- err.Error()
		return nil, err
	}

	trains, err := a.GenerateTrainsList(route, args)
	if err != nil {
		a.LogChan <- err.Error()
		return nil, err
	}

	return trains, nil
}

func (a *App) GenerateTrainsList(route entity.Route, args entity.RouteArgs) ([]entity.Train, error) {
	trainsAnswer := []entity.Train{}

	trains, err := getTrainsList(route, args)
	if err != nil {
		return trainsAnswer, err
	}

	for _, val := range trains {
		data, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		formatKey := fmt.Sprintf("%s_%s", val.Number, val.Date0)

		err = a.Cache.Set(formatKey, data)
		if err != nil {
			return nil, err
		}
		val.ID = formatKey
		trainsAnswer = append(trainsAnswer, val)
	}

	return trainsAnswer, nil
}

func (a *App) GetStationCodes(target, source string) (int, int, error) {
	var code = make(chan GoroutineAnswer)
	var answers = map[string]int{}
	go func() {
		data, err := a.Routes.GetDirectionsCode(target)
		if err != nil {
			a.LogChan <- err.Error()
		}
		answer := GoroutineAnswer{
			Code:    data,
			Station: "target",
		}
		code <- answer
	}()
	go func() {
		data, err := a.Routes.GetDirectionsCode(source)
		if err != nil {
			a.LogChan <- err.Error()
		}
		answer := GoroutineAnswer{
			Code:    data,
			Station: "source",
		}
		code <- answer
	}()
	for {
		select {
		case val := <-code:
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

func (a *App) SaveInfoAboutTrain(trainID string) (string, error) {
	train := entity.Train{}

	data, err := a.Cache.Get(trainID)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &train)
	if err != nil {
		return "", err
	}

	trainID, err = a.Trains.Create(train)
	if err != nil {
		return "", err
	}

	return trainID, nil
}

func (a *App) Run(refreshTimeSec string) {
	minutes, _ := time.ParseDuration(refreshTimeSec)
	ticker := time.NewTicker(minutes)
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				return
			}

			trains, err := a.Trains.ReadMany()
			if err != nil {
				log.Printf("%s\n", err)
			}

			for _, val := range trains {
				if a.CheckAndRefreshTrainInfo(val) {
					err := a.Trains.Update(val)
					if err != nil {
						a.LogChan <- fmt.Sprintf("%s", err.Error())
						continue
					}
				}
			}
		}
	}
}

func (a *App) CheckAndRefreshTrainInfo(train entity.Train) bool {
	//RID
	rid, cookies, err := a.Routes.GetRid(trainToArgs(train))

	if err != nil {
		a.LogChan <- fmt.Sprintf("App->CheckAndRefreshTrainInfo: Error while requesting rid from RZD API%s", err.Error())
		return false
	}
	train.QueryArgs.Rid = strconv.FormatInt(rid.RID, 10)

	newRoute, err := a.Routes.GetInfoAboutOneTrain(train, cookies)
	if err != nil {
		a.LogChan <- fmt.Sprintf("App->CheckAndRefreshTrainInfo: Error while handling answer about one train %s", err.Error())
		return false
	}

	trains, err := getTrainsList(newRoute, train.QueryArgs)
	if err != nil {
		a.LogChan <- fmt.Sprintf("App->GenerateTrainsList: Got empty route array")
		return false
	}
	for _, val := range trains {
		if a.GetDiff(train, val) {
			err := a.Trains.Update(val)
			if err != nil {
				a.LogChan <- fmt.Sprintf("App->CheckAndRefreshTrainInfo: Error while update train in db %s", err.Error())
				return false
			}
		} else {
			a.LogChan <- "App->CheckAndRefreshTrainInfo: don't find diff"
		}
	}

	return true
}

func (a *App) GetDiff(oldTrain entity.Train, newTrain entity.Train) bool {
	for key, val := range oldTrain.Seats {
		if newTrain.Seats[key].SeatsCount > val.SeatsCount && newTrain.Seats[key].Price == val.Price {
			if val.Chosen == true {
				return true
			}
		}
	}
	return false
}

func getTrainsList(route entity.Route, args entity.RouteArgs) ([]entity.Train, error) {
	trains := []entity.Train{}

	if len(route.Tp) == 0 {
		return nil, errors.New(fmt.Sprintf("App->GenerateTrainsList: Got empty route array"))
	}

	for _, val := range route.Tp[0].List {
		var seats = entity.Seats{
			entity.CSeatsType: {
				SeatsCount: 0,
				Price:      "0",
				Chosen:     false,
			},
			entity.SSeatsType: {
				SeatsCount: 0,
				Price:      "0",
				Chosen:     false,
			},
			entity.SVSeatsType: {
				SeatsCount: 0,
				Price:      "0",
				Chosen:     false,
			},
			entity.PSeatsType: {
				SeatsCount: 0,
				Price:      "0",
				Chosen:     false,
			},
		}

		for _, seatsInfo := range val.ServiceCategories {
			seats[fromRzdToCore(seatsInfo.TypeLoc)] = entity.Seat{
				SeatsCount: seatsInfo.FreeSeats,
				Price:      seatsInfo.Price,
			}
		}

		trains = append(trains, entity.Train{
			Number:    val.Number,
			Route0:    val.Route0,
			Route1:    val.Route1,
			TrDate0:   val.TrDate0,
			TrTime0:   val.TrTime0,
			Station:   val.Station,
			Station1:  val.Station1,
			Date0:     val.Date0,
			Time0:     val.Time0,
			Date1:     val.Date1,
			Time1:     val.Time1,
			Seats:     seats,
			QueryArgs: args,
		})
	}

	return trains, nil
}

func (a *App) AddUser(user *entity.User) (bool, error) {
	ok, err := a.Users.Create(user)
	if err != nil {
		return ok, err
	}
	return ok, nil
}

//???? nahyia
func (a *App) UpdateUserTrainInfo(user *entity.User) error {
	err := a.Users.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteUser(user *entity.User) error {
	err := a.Users.Delete(user)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) GetUsersList() ([]entity.User, error) {
	users, err := a.Users.ReadMany()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (a *App) UsersCount() (int, error) {
	users, err := a.Users.ReadMany()
	if err != nil {
		return 0, err
	}
	return len(users), nil
}

func (a *App) SaveTrainInUser(userID string, trainID string) error {
	user := &entity.User{
		UserTelegramID: userID,
	}
	savedUser, err := a.Users.ReadOne(user)
	if err != nil {
		return err
	}

	for _, val := range savedUser.TrainIDS {
		if val == trainID {
			return nil
		}
	}
	savedUser.TrainIDS = append(savedUser.TrainIDS, trainID)
	err = a.Users.Update(savedUser)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) CheckUsers(start, end int64) ([]*entity.User, error) {
	users, err := a.Users.ReadSection(start, end)
	notifyUsers := []*entity.User{}

	if err != nil {
		return nil, err
	}
	for _, val := range users {
		trains := val.TrainIDS
		for i := range trains {
			train, err := a.Trains.ReadOne(trains[i])
			if err != nil {
				a.LogChan <- err.Error()
				a.LogChan <- fmt.Sprintf("App->GenerateTrainsList: Can't get train")
			}
			if a.CheckAndRefreshTrainInfo(train) {
				notifyUsers = append(notifyUsers, val)
				a.LogChan <- fmt.Sprintf("App->GenerateTrainsList: all good in train - %v", train)
			}
		}
	}
	return notifyUsers, nil
}

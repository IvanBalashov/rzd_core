package usecase

import (
	"rzd/app/gateways/cache_gateway"
	"rzd/app/gateways/rzd_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
)

// TODO: Think about how correct work with error messages.
type App struct {
	Trains  trains_gateway.TrainsGateway
	Users   users_gateway.UsersGateway
	Routes  rzd_gateway.RzdGateway
	Cache   cache_gateway.CacheGateway
	LogChan chan string
	Rid     string
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

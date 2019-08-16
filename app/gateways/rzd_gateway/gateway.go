package rzd_gateway

import (
	"net/http"
	"rzd/app/entity"
)

//https://github.com/visavi/rzd-api
type RzdGateway interface {
	GetRoutes(args *entity.RouteArgs, cookie []*http.Cookie) (*entity.Route, error)
	GetDirectionsCode(source string) (int, error)
	GetRid(args *entity.RidArgs) (*entity.Rid, []*http.Cookie, error)
	GetInfoAboutOneTrain(train *entity.Train, cookie []*http.Cookie) (*entity.Route, error)
}

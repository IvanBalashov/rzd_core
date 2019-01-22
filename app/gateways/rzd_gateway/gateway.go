package rzd_gateway

import "rzd/app/entity"

//https://github.com/visavi/rzd-api
type RzdGateway interface {
	GetRoutes(args entity.RouteArgs) (entity.Route, error)
	GetDirectionsCode(source string) (int, error)
	GetRid(args entity.RidArgs) (entity.Rid, error)
}

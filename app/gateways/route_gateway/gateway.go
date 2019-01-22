package route_gateway

import "rzd/app/entity"

//https://github.com/visavi/rzd-api
type RouteGateway interface {
	GetRoutes(args entity.RouteArgs) (entity.Route, error)
}

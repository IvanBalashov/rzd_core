package middleware

import "rzd/app/usecase"

type EventLayer struct {
	App usecase.Usecase
}

func InitMiddleWares(app usecase.Usecase) EventLayer {
	return EventLayer{
		App: app,
	}
}

package middleware

import "rzd/app/usecase"

type AppLayer struct {
	App usecase.Usecase
}

func NewEventLayer(app usecase.Usecase) AppLayer {
	return AppLayer{
		App: app,
	}
}

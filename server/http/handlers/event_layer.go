package handlers

import "rzd/app/usecase"

type EventLayer struct {
	App usecase.Usecase
}

func NewEventLayer(app usecase.Usecase) *EventLayer {
	return &EventLayer{
		App: app,
	}
}

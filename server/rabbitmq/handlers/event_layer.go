package handlers

import "rzd/app/usecase"

type EventLayer struct {
	App       usecase.Usecase
	LogChanel chan string
}

func NewEventLayer(app usecase.Usecase, logChan chan string) *EventLayer {
	return &EventLayer{
		App:       app,
		LogChanel: logChan,
	}
}

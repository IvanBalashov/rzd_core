package usecase

import (
	"log"
	"time"
)

func (a *App) Run(refreshTimeSec string) {
	minutes, _ := time.ParseDuration(refreshTimeSec)
	ticker := time.NewTicker(minutes)
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				return
			}
			_, err := a.Trains.ReadMany()
			if err != nil {
				log.Printf("%s\n", err)
			}

		}
	}
}

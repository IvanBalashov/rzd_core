package usecase

import (
	"fmt"
	"log"
	"time"
)

func (a *App) Run(refreshTimeSec string) {
	minutes, _ := time.ParseDuration(fmt.Sprintf("%ss", refreshTimeSec))
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

		}
	}
}

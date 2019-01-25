package usecase

import (
	"fmt"
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
				panic(err)
			}

			for key, val := range trains {
				fmt.Printf("%s - %s\n", key, val)
			}
		}
	}
}

func checkTrain() {

}

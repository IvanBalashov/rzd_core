package usecase

type GoroutineAnswer struct {
	Code    int
	Station string
}

func (a *App) GetStationCodes(target, source string) (int, int, error) {
	var code = make(chan GoroutineAnswer)
	var answers = map[string]int{}
	go func() {
		data, err := a.Routes.GetDirectionsCode(target)
		if err != nil {
			a.LogChan <- err.Error()
		}
		answer := GoroutineAnswer{
			Code:    data,
			Station: "target",
		}
		code <- answer
	}()
	go func() {
		data, err := a.Routes.GetDirectionsCode(source)
		if err != nil {
			a.LogChan <- err.Error()
		}
		answer := GoroutineAnswer{
			Code:    data,
			Station: "source",
		}
		code <- answer
	}()
	for {
		select {
		case val := <-code:
			if val.Station == "target" {
				answers["target"] = val.Code
			} else {
				answers["source"] = val.Code
			}
		}
		if len(answers) == 2 {
			break
		}
	}
	return answers["target"], answers["source"], nil
}

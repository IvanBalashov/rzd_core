package entity

/*
	dir - 0 только в один конец, 1 - туда-обратно
	tfl - тип поезда (1- все, 2 - дальнего следования, 3- электрички)
	checkSeats - 1, 0 - поиск в поездах только если есть свободные места
	code0 - код станции отправления
	code1 - код станции прибытия
	dt0 - дата отправления
*/

type RidArgs struct {
	Dir          string `json:"dir"`
	Tfl          string `json:"tfl"`
	CheckSeats   string `json:"checkSeats"`
	Code0        string `json:"code0"`
	Code1        string `json:"code1"`
	Dt0          string `json:"dt0"`
	WithOutSeats string `json:"withoutSeats"`
	Version      string `json:"version"`
}

func (r *RidArgs) ToMap() map[string]string {
	return map[string]string{
		"dir":          r.Dir,
		"tfl":          r.Tfl,
		"code0":        r.Code0,
		"code1":        r.Code1,
		"dt0":          r.Dt0,
		"checkSeats":   r.CheckSeats,
		"withoutSeats": r.WithOutSeats,
		"version":      r.Version,
	}
}

package entity

type Train struct {
	ID       string             `json:"id"`
	Number   string             `json:"number"`
	Route0   string             `json:"route_0"`
	Route1   string             `json:"route_1"`
	TrDate0  string             `json:"tr_date_0"`
	TrTime0  string             `json:"tr_time_0"`
	Station  string             `json:"station"`
	Station1 string             `json:"station_1"`
	Date0    string             `json:"date_0"`
	Time0    string             `json:"time_0"`
	Date1    string             `json:"date_1"`
	Time1    string             `json:"time_1"`
	Seats    map[SeatsType]Seat `json:"seats"`
	// But we still need generate new rid...
	QueryArgs RouteArgs `json:"query_args"`
}

type Seat struct {
	SeatsCount int  `json:"seats_count"`
	Price      int  `json:"price"`
	Chosen     bool `json:"chosen"`
}

type SeatsType string

const (
	CSeatsType  SeatsType = "Купе"
	SSeatsType  SeatsType = "Сидячий"
	SVSeatsType SeatsType = "СВ"
	PSeatsType  SeatsType = "Плацкартный"
)

type Seats = map[SeatsType]Seat

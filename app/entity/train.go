package entity

type Train struct {
	ID       string
	//Type     string // TODO: think this don't needed
	Number   string
	//Brand    string // TODO: think this don't needed
	Route0   string // TODO: think this don't needed
	Route1   string // TODO: think this don't needed
	TrDate0  string
	TrTime0  string
	Station  string
	Station1 string
	Date0    string // TODO: think this don't needed
	Time0    string
	Date1    string // TODO: think this don't needed
	Time1    string // TODO: think this don't needed
	Seats    map[SeatsType]Seat
	// But we still need generate new rid...
	QueryArgs RouteArgs
}

type Seat struct {
	SeatsCount int
	Price      string
	Chosen     bool
}

type SeatsType string

const (
	CSeatsType  SeatsType = "Купе"
	SSeatsType  SeatsType = "Сидячий"
	SVSeatsType SeatsType = "СВ"
	PSeatsType  SeatsType = "Плацкартный"
)

type Seats = map[SeatsType]Seat

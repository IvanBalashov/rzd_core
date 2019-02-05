package entity

type Train struct {
	ID       string
	Type     string
	Number   string
	Brand    string
	Route0   string
	Route1   string
	TrDate0  string
	TrTime0  string
	Station  string
	Station1 string
	Date0    string
	Time0    string
	Date1    string
	Time1    string
	Seats    []Seats
	// But we still need generate new rid...
	QueryArgs RouteArgs
}

type Seats struct {
	SeatsName  string
	SeatsCount int
	Price      int
}

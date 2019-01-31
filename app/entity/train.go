package entity

// TODO: think this shit don't needed.
type Train struct {
	ID       string
	Code0    string
	Code1    string
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
}

type Seats struct {
	SeatsName  string
	SeatsCount int
	Price      int
}

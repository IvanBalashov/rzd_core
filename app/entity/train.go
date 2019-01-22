package entity

type Train struct {
	Type       string
	Number     string
	Brand      string
	Route0     string
	Route1     string
	TrDate0    string
	TrTime0    string
	Station    string
	Station1   string
	Date0      string
	Time0      string
	Date1      string
	Time1      string
	Class      string
	SeatsCount string
	Price      string
}

//number, type, brand, route0, route1, trTime0, station, station1, date0, time0, date1, time1, class, seatsCount, price
func (t *Train) GetArgs() (string, string, string, string, string, string, string, string, string, string, string, string, string) {
	return t.Type, t.Number, t.Brand, t.Route0, t.Route1, t.TrDate0, t.TrTime0, t.Station, t.Station1, t.Date0, t.Date1, t.Time1, t.SeatsCount
}

/* CREATE TABLE trains (
id SERIAL PRIMARY KEY,
number varchar(20),
type varchar(20),
brand varchar(20),
route0 varchar(100),
route1 varchar(100),
trTime0 varchar(10),
station varchar(100),
station1 varchar(100),
date0 varchar(50),
time0 varchar(50),
date1 varchar(50),
time1 varchar(50),
class varchar(50),
seatsCount varchar(20),
price varchar(20));
*/

package entity

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
func (t *Train) GetArgs() []string {
	fields := []string{}
	fields = append(fields, t.Number)
	fields = append(fields, t.Type)
	fields = append(fields, t.Route0)
	fields = append(fields, t.Route1)
	fields = append(fields, t.TrDate0)
	fields = append(fields, t.TrTime0)
	fields = append(fields, t.Station)
	fields = append(fields, t.Station1)
	fields = append(fields, t.Date0)
	fields = append(fields, t.Time0)
	fields = append(fields, t.Class)
	fields = append(fields, t.SeatsCount)
	fields = append(fields, t.Price)
	return fields
}

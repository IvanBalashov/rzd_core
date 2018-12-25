package entity

type Route struct {
	Result string `json:"result"`
	Tp     []Tp   `json:"tp"`
}
type Tp struct {
	From        string `json:"from"`
	FromCode    uint64 `json:"fromCode"`
	Where       string `json:"where"`
	WhereCode   uint64 `json:"whereCode"`
	Date        string `json:"date"`
	NoSeats     bool   `json:"noSeats"`
	DefShowTime string `json:"defShowTime"`
	State       string `json:"state"`
	List        []List `json:"list"`
}

type List struct {
	Number            string              `json:"number"`
	Type              int                 `json:"type"`
	Brand             string              `json:"brand"`
	Route0            string              `json:"route0"`
	Route1            string              `json:"route1"`
	TrDate0           string              `json:"trDate0"`
	TrTime0           string              `json:"trTime0"`
	Station           string              `json:"station0"`
	Station1          string              `json:"station1"`
	Date0             string              `json:"date0"`
	Time0             string              `json:"time0"`
	Date1             string              `json:"date1"`
	Time1             string              `json:"time1"`
	ServiceCategories []ServiceCategories `json:"cars"`
}

type ServiceCategories struct {
	TypeLoc   string `json:"typeLoc"`
	FreeSeats int    `json:"freeSeats"`
	Price     int    `json:"tarif"`
}

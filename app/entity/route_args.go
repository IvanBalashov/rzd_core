package entity

type RouteArgs struct {
	Dir          string `json:"dir"`
	Tfl          string `json:"tfl"`
	CheckSeats   string `json:"checkSeats"`
	Code0        string `json:"code0"`
	Code1        string `json:"code1"`
	Dt0          string `json:"dt0"`
	Dt1          string `json:"dt_1"`
	WithOutSeats string `json:"withoutSeats"`
	Version      string `json:"version"`
	Rid          string `json:"rid,omitempty"`
}

func (r *RouteArgs) ToMap() map[string]string {
	return map[string]string{
		"dir":        r.Dir,
		"tfl":        r.Tfl,
		"checkSeats": r.CheckSeats,
		"code0":      r.Code0,
		"dt0":        r.Dt0,
		"code1":      r.Code1,
		"dt1":        r.Dt1,
		"rid":        r.Rid,
	}
}

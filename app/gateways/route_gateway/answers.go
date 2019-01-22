package route_gateway

type Rid struct {
	RID       int64  `json:"RID"`
	Result    string `json:"result"`
	Timestamp string `json:"timestamp"`
}

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

type Codes struct {
	Name string `json:"n"`
	Code int    `json:"c"`
}

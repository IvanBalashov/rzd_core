package route_gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"rzd/app/entity"
	"strconv"
	"time"
)

// TODO: Create self APIClient for single user????
type APIClient struct {
	APIUrl string //http://pass.rzd.ru/timetable/public/ru
	Code1  int    //?layer_id=5827
	Code2  int    //?layer_id=5764
	Code3  int    //?layer_id=5804
}

func NewRestAPIClient(url string, code1, code2, code3 int) APIClient {
	return APIClient{
		APIUrl: url,
		Code1:  code1,
		Code2:  code2,
		Code3:  code3,
	}
}

func (a *APIClient) GetRoutes(args entity.RouteArgs) entity.Route {
	rid := a.getRid(RidArgs{
		Dir:          args.Dir,
		Tfl:          args.Tfl,
		CheckSeats:   args.CheckSeats,
		Code0:        args.Code0,
		Code1:        args.Code1,
		Dt0:          args.Dt0,
		WithOutSeats: args.WithOutSeats,
		Version:      args.Version,
	})

	args.Rid = strconv.FormatInt(rid.RID, 10)
	// this sleep needed coz server need time to save rid.
	time.Sleep(time.Second * 1)

	route := entity.Route{}

	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetFormData(args.ToMap()).
		SetQueryParam("layer_id", "5827").
		Post(a.APIUrl)
	if err != nil {
		panic(err)
	}

	body := resp.Body()
	str := bytes.NewBuffer(body).String()
	fmt.Printf("%s\n", str)

	err = json.Unmarshal(resp.Body(), &route)
	if err != nil {
		panic(err)
	}

	return route
}

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

func (a *APIClient) getRid(args RidArgs) Rid {
	rid := Rid{}

	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetFormData(args.ToMap()).
		SetQueryParam("layer_id", "5827").
		Post(a.APIUrl)
	if err != nil {
		panic(err)
	}

	body := resp.Body()
	str := bytes.NewBuffer(body).String()
	fmt.Printf("%s\n", str)
	err = json.Unmarshal(body[10:], &rid)
	if err != nil {
		panic(err)
	}

	return rid
}

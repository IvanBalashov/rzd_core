package route_gateway

import (
	"encoding/json"
	"gopkg.in/resty.v1"
	"log"
	"rzd/app/entity"
	"strconv"
	"strings"
	"time"
)

// TODO: Create self APIClient for single user????
// FIXME: Refactor variables name!!!
// FIXME: rewrite all panic methods.
type APIClient struct {
	// OK, im hope what this url don't changes.
	PassRzdUrl string //http://pass.rzd.ru/timetable/public/ru
	RzdUrl     string //http://www.rzd.ru/
	Code1      int    //?layer_id=5827
	Code2      int    //?layer_id=5764
	Code3      int    //?layer_id=5804
}

func NewRestAPIClient(url string, code1, code2, code3 int) APIClient {
	return APIClient{
		PassRzdUrl: url,
		Code1:      code1,
		Code2:      code2,
		Code3:      code3,
	}
}

func (a *APIClient) GetRoutes(args entity.RouteArgs) (entity.Route, error) {
	// im think is good practise to provide inner errors in app layer.
	rid, err := a.getRid(RidArgs{
		Dir:          args.Dir,
		Tfl:          args.Tfl,
		CheckSeats:   args.CheckSeats,
		Code0:        args.Code0,
		Code1:        args.Code1,
		Dt0:          args.Dt0,
		WithOutSeats: args.WithOutSeats,
		Version:      args.Version,
	})
	if err != nil {
		return entity.Route{}, err
	}

	args.Rid = strconv.FormatInt(rid.RID, 10)
	// this sleep needed coz server need time to save rid.
	// FIXME: lower w8ing time.
	time.Sleep(time.Second)

	route := entity.Route{}

	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetFormData(args.ToMap()).
		SetQueryParam("layer_id", "5827").
		Post(a.PassRzdUrl)
	if err != nil {
		log.Printf("Gateways->Rzd_Gateway->GetRoutes: Error in request to RZD Api - %s\n", err)
		return entity.Route{}, err
	}

	err = json.Unmarshal(resp.Body(), &route)
	if err != nil {
		log.Printf("Gateways->Rzd_Gateway->GetRoutes: Error in unmarshal anwer from RZD Api - %s\n", err)
		return entity.Route{}, err
	}

	return route, nil
}

func (a *APIClient) getRid(args RidArgs) (Rid, error) {
	rid := Rid{}

	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetFormData(args.ToMap()).
		SetQueryParam("layer_id", "5827").
		Post(a.PassRzdUrl)
	if err != nil {
		log.Printf("Gateways->Rzd_Gateway->getRid: Error in request to RZD Api - %s\n", err)
		return Rid{}, err
	}

	body := resp.Body()
	// need clear first 10 symbols coz answer from rzd api have "\n" 5 symbols to move cursor down.
	err = json.Unmarshal(body[10:], &rid)
	if err != nil {
		log.Printf("Gateways->Rzd_Gateway->getRid: Error in unmarshal anwer from RZD Api - %s\n", err)
		return Rid{}, err
	}

	return rid, nil
}

//Coz all rzd rest api distributed on two entry points - pass.rzd.ru and rzd.ru.
func (a *APIClient) GetDirectionsCode(source string) (int, error) {
	answer := []Codes{}
	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetQueryParam("", "").
		Get(a.RzdUrl)
	if err != nil {
		log.Printf("Gateways->Rzd_Gateway->GetDirectionsCode: Error in request to RZD Api - %s\n", err)
		return 0, err
	}

	err = json.Unmarshal(resp.Body(), &answer)
	if err != nil {
		log.Printf("Gateways->Rzd_Gateway->GetDirectionsCode: Error in unmarshal anwer from RZD Api - %s\n", err)
		return 0, err
	}
	for i := range answer {
		if strings.ToLower(answer[i].Name) == strings.ToLower(source) {
			return answer[i].Code, nil
		}
	}
	return 0, nil
}

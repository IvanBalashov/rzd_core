package rzd_gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/resty.v1"
	"net/http"
	"rzd/app/entity"
	"strings"
)

// TODO: Create self APIClient for single user????
// FIXME: Refactor variables name!!!
// Codes can be hardcoded.?
type APIClient struct {
	// OK, im hope what this url don't changes.
	PassRzdUrl string //http://pass.rzd.ru/timetable/public/ru
	RzdUrl     string //http://www.rzd.ru/
	Code1      int    //?layer_id=5827
	Code2      int    //?layer_id=5764
	Code3      int    //?layer_id=5804
}

func NewRestAPIClient(passUrl, rzdUrl string, code1, code2, code3 int) APIClient {
	return APIClient{
		PassRzdUrl: passUrl,
		RzdUrl:     rzdUrl,
		Code1:      code1,
		Code2:      code2,
		Code3:      code3,
	}
}

func (a *APIClient) GetRoutes(args entity.RouteArgs, cookies []*http.Cookie) (entity.Route, error) {
	route := entity.Route{}

	for key := range cookies {
		resty.SetCookie(cookies[key])
	}

	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetFormData(args.ToMap()).
		SetQueryParam("layer_id", "5827").
		Post(a.PassRzdUrl)
	if err != nil {
		return entity.Route{},
			errors.New(fmt.Sprintf("Gateways->Rzd_Gateway->GetRoutes: Error in request to RZD Api - %s", err))
	}

	err = json.Unmarshal(resp.Body(), &route)
	if err != nil {
		return entity.Route{},
			errors.New(fmt.Sprintf("Gateways->Rzd_Gateway->GetRoutes: Error in unmarshal anwer from RZD Api - %s", err))
	}

	return route, nil
}

func (a *APIClient) GetRid(args entity.RidArgs) (entity.Rid, []*http.Cookie, error) {
	rid := entity.Rid{}

	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetFormData(args.ToMap()).
		SetQueryParam("layer_id", "5827").
		Post(a.PassRzdUrl)
	if err != nil {
		return entity.Rid{}, nil,
			errors.New(fmt.Sprintf("Gateways->Rzd_Gateway->getRid: Error in request to RZD Api - %s", err))
	}

	cookies := resp.Cookies()

	body := resp.Body()
	err = json.Unmarshal(body, &rid)
	if err != nil {
		return entity.Rid{}, nil,
			errors.New(fmt.Sprintf("Gateways->Rzd_Gateway->getRid: Error in unmarshal anwer from RZD Api - %s\n", err))
	}

	return rid, cookies, nil
}

//Coz all rzd rest api distributed on two entry points - pass.rzd.ru and rzd.ru.
func (a *APIClient) GetDirectionsCode(source string) (int, error) {
	answer := []Codes{}
	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetQueryParam("stationNamePart", strings.ToUpper(source[:4])).
		SetQueryParam("lang", "ru").
		SetQueryParam("lat", "0").
		SetQueryParam("compactMode", "y").
		Get(a.RzdUrl)
	if err != nil {
		return 0,
			errors.New(fmt.Sprintf("Gateways->Rzd_Gateway->GetDirectionsCode: Error in request to RZD Api - %s", err))
	}

	err = json.Unmarshal(resp.Body(), &answer)
	if err != nil {
		return 0,
			errors.New(fmt.Sprintf("Gateways->Rzd_Gateway->GetDirectionsCode: Error in unmarshal anwer from RZD Api - %s", err))
	}
	for i := range answer {
		if strings.Contains(strings.ToLower(answer[i].Name), strings.ToLower(source)) {
			return answer[i].Code, nil
		}
	}
	return 0, nil
}

//FIXME: NOT TESTED METHOD
func (a *APIClient) GetInfoAboutOneTrain(train entity.Train, cookies []*http.Cookie) (entity.Route, error) {
	answer := entity.Route{}

	for key := range cookies {
		resty.SetCookie(cookies[key])
	}

	resp, err := resty.R().
		SetHeader("Accept", "application/json").
		SetQueryParam("layer_id", "5827").
		SetQueryParam("tnum0", train.Number).
		SetQueryParams(train.QueryArgs.ToMap()).
		Get(a.PassRzdUrl)

	// FIXME
	err = json.Unmarshal(resp.Body(), &answer)
	if err != nil {
		return entity.Route{},
			errors.New(fmt.Sprintf("Gateways->Rzd_Gateway->GetInfoAboutOneTrain: Error in request to RZD Api - %s", err))
	}

	return entity.Route{}, nil
}

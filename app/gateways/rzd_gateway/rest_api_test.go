package rzd_gateway

import (
	"encoding/json"
	"gopkg.in/go-playground/assert.v1"
	"io/ioutil"
	"net/http"
	"rzd/app/entity"
	"testing"
)

var api = APIClient{
	"http://localhost:9090/timetable/public/ru",
	"http://localhost:9090/suggester/station",
	5827,
	5764,
	5804,
}

func TestNewRestAPIClient(t *testing.T) {
	newApi := NewRestAPIClient(
		"http://localhost:9090/timetable/public/ru",
		"http://localhost:9090/suggester/station",
		5827,
		5764,
		5804,
	)

	assert.Equal(t, api, newApi)
}

func TestAPIClient_GetDirectionsCode(t *testing.T) {
	expectCode := 2000000

	code, err := api.GetDirectionsCode("Москва")
	if err != nil {
		t.Errorf("error in GetDirectionsCode - %s\n", err)
		t.FailNow()
	}

	assert.Equal(t, expectCode, code)
}

func TestAPIClient_GetRid(t *testing.T) {
	expectedRid := entity.Rid{}
	expectedCookies := []*http.Cookie{}
	testingArgs := entity.RidArgs{
		Dir:          "0",
		Tfl:          "1",
		Code0:        "2000000",
		Code1:        "2010000",
		Dt0:          "8.02.2019",
		CheckSeats:   "0",
		WithOutSeats: "y",
		Version:      "v.2018",
	}

	data, err := ioutil.ReadFile("./datasets/expected_rid.json")
	if err != nil {
		t.Errorf("error in can't read expected data - %s\n", err)
		t.FailNow()
	}

	err = json.Unmarshal(data, &expectedRid)
	if err != nil {
		t.Errorf("error in can't read expected data - %s\n", err)
		t.FailNow()
	}

	rid, cookies, err := api.GetRid(testingArgs)
	if err != nil {
		t.Errorf("error in GetRid - %s\n", err)
		t.FailNow()
	}

	assert.Equal(t, expectedRid, rid)
	assert.Equal(t, expectedCookies, cookies)
}

func TestAPIClient_GetRoutes(t *testing.T) {
	testingArgs := entity.RouteArgs{
		Dir:          "0",
		Tfl:          "1",
		Code0:        "2000000",
		Code1:        "2010000",
		Dt0:          "8.02.2019",
		CheckSeats:   "0",
		WithOutSeats: "y",
		Version:      "v.2018",
		Rid:          "6862063710",
	}
	testingCookies := []*http.Cookie{}
	data, err := ioutil.ReadFile("./datasets/expected_route.json")
	if err != nil {
		t.Errorf("error in can't read expected data - %s\n", err)
		t.FailNow()
	}

	expectedRoute := entity.Route{}
	err = json.Unmarshal(data, &expectedRoute)
	if err != nil {
		t.Errorf("error in GetRoutes - %s\n", err)
		t.FailNow()
	}

	route, err := api.GetRoutes(testingArgs, testingCookies)
	if err != nil {
		t.Errorf("error in GetRoutes - %s\n", err)
		t.FailNow()
	}

	assert.Equal(t, expectedRoute, route)

}

package trains_gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"
	"os"
	"rzd/app/entity"
	"testing"
	"time"
)

func SetupConnection() (*MongoTrains, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	connUrl := os.Getenv("MONGODB_URL")
	if connUrl == "" {
		return nil, errors.New("Can't connect to mongo_db. MONGODB_URL not setted.")
	}

	client, err := mongo.Connect(ctx, connUrl)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	col := client.Database("rzd").Collection("trains")

	return &MongoTrains{
		CLI:    client,
		Trains: col,
	}, nil
}

type DeletedTrain struct {
	ID string `json:"id"`
}

func cleanFunc(t *testing.T, clean bool, client *mongo.Collection, obj interface{}) {
	if clean {
		res, err := client.DeleteMany(context.Background(), obj)
		if err != nil {
			t.Fatal(err)
		}
		if res.DeletedCount == 0 {
			t.Fatal(errors.New(fmt.Sprintf("can't delete %v", obj)))
		}
	}
}

func TestMongoTrains_Create(t *testing.T) {
	trains, err := SetupConnection()
	if err != nil {
		t.Fatal(err)
	}
	a := assert.New(t)
	testCases := []struct {
		name       string
		train      *entity.Train
		exceptedID string
		cleanFlag  bool
		checkErr   func(err error, msgAndArgs ...interface{}) bool
	}{
		{
			"success",
			&entity.Train{
				ID:       "test",
				Number:   "123",
				Route0:   "msk",
				Route1:   "spb",
				TrDate0:  "01.01.1970",
				TrTime0:  "01.01.1970",
				Station:  "main",
				Station1: "main2",
				Date0:    "01.01.1970",
				Time0:    "00:00",
				Date1:    "01.01.1970",
				Time1:    "12:00",
				Seats: map[entity.SeatsType]entity.Seat{
					entity.CSeatsType: {
						0, 0, false,
					},
					entity.SVSeatsType: {
						0, 0, false,
					},
					entity.SSeatsType: {
						0, 0, false,
					},
					entity.PSeatsType: {
						0, 0, false,
					},
				},
				QueryArgs: entity.RouteArgs{},
			},
			"123123",
			true,
			a.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hex, err := trains.Create(tc.train)
			tc.checkErr(err)
			a.NotEqual(hex, "")

			cleanFunc(t, tc.cleanFlag, trains.Trains, &DeletedTrain{ID: tc.train.ID})
		})
	}
}

//TODO: 1) Добавить контескст для всех запросов, начиная от хендлеров
//TODO: 2) Сделать внутниние структуры для работы с монгой, а то чет какая то хуйня получается пока что

func TestMongoTrains_Delete(t *testing.T) {
	trains, err := SetupConnection()
	if err != nil {
		t.Fatal(err)
	}
	a := assert.New(t)
	testCases := []struct {
		name      string
		train     *entity.Train
		cleanFlag bool
		checkErr  func(err error, msgAndArgs ...interface{}) bool
	}{
		{
			"success",
			&entity.Train{
				ID:       "test",
				Number:   "123",
				Route0:   "msk",
				Route1:   "spb",
				TrDate0:  "01.01.1970",
				TrTime0:  "01.01.1970",
				Station:  "main",
				Station1: "main2",
				Date0:    "01.01.1970",
				Time0:    "00:00",
				Date1:    "01.01.1970",
				Time1:    "12:00",
				Seats: map[entity.SeatsType]entity.Seat{
					entity.CSeatsType: {
						0, 0, false,
					},
					entity.SVSeatsType: {
						0, 0, false,
					},
					entity.SSeatsType: {
						0, 0, false,
					},
					entity.PSeatsType: {
						0, 0, false,
					},
				},
				QueryArgs: entity.RouteArgs{},
			},
			false,
			a.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := trains.Trains.InsertOne(context.Background(), tc.train)
			a.NoError(err)

			err = trains.Delete(tc.train)
			tc.checkErr(err)
		})
	}
}

func TestMongoTrains_ReadOne(t *testing.T) {
	trains, err := SetupConnection()
	if err != nil {
		t.Fatal(err)
	}
	a := assert.New(t)
	testCases := []struct {
		name          string
		expectedTrain *entity.Train
		cleanFlag     bool
		checkErr      func(err error, msgAndArgs ...interface{}) bool
	}{
		{
			"success",
			&entity.Train{
				ID:       "test",
				Number:   "123",
				Route0:   "msk",
				Route1:   "spb",
				TrDate0:  "01.01.1970",
				TrTime0:  "01.01.1970",
				Station:  "main",
				Station1: "main2",
				Date0:    "01.01.1970",
				Time0:    "00:00",
				Date1:    "01.01.1970",
				Time1:    "12:00",
				Seats: map[entity.SeatsType]entity.Seat{
					entity.CSeatsType: {
						0, 0, false,
					},
					entity.SVSeatsType: {
						0, 0, false,
					},
					entity.SSeatsType: {
						0, 0, false,
					},
					entity.PSeatsType: {
						0, 0, false,
					},
				},
				QueryArgs: entity.RouteArgs{},
			},
			true,
			a.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := trains.Trains.InsertOne(context.Background(), tc.expectedTrain)
			a.NoError(err)

			actualTrain, err := trains.ReadOne(tc.expectedTrain)
			tc.checkErr(err)

			assert.Equal(t, tc.expectedTrain, actualTrain)

			cleanFunc(t, tc.cleanFlag, trains.Trains, &DeletedTrain{ID: tc.expectedTrain.ID})
		})
	}
}

func TestMongoTrains_Update(t *testing.T) {
	trains, err := SetupConnection()
	if err != nil {
		t.Fatal(err)
	}
	a := assert.New(t)
	testCases := []struct {
		name          string
		expectedTrain *entity.Train
		cleanFlag     bool
		checkErr      func(err error, msgAndArgs ...interface{}) bool
	}{
		{
			"success",
			&entity.Train{
				ID:       "test",
				Number:   "123",
				Route0:   "msk",
				Route1:   "spb",
				TrDate0:  "01.01.1970",
				TrTime0:  "01.01.1970",
				Station:  "main",
				Station1: "main2",
				Date0:    "01.01.1970",
				Time0:    "00:00",
				Date1:    "01.01.1970",
				Time1:    "12:00",
				Seats: map[entity.SeatsType]entity.Seat{
					entity.CSeatsType: {
						0, 0, false,
					},
					entity.SVSeatsType: {
						0, 0, false,
					},
					entity.SSeatsType: {
						0, 0, false,
					},
					entity.PSeatsType: {
						0, 0, false,
					},
				},
				QueryArgs: entity.RouteArgs{},
			},
			true,
			a.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := trains.Trains.InsertOne(context.Background(), tc.expectedTrain)
			a.NoError(err)

			err = trains.Update(tc.expectedTrain)
			tc.checkErr(err)

			cleanFunc(t, tc.cleanFlag, trains.Trains, &DeletedTrain{ID: tc.expectedTrain.ID})
		})
	}
}

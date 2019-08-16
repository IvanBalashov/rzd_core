package users_gateway

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

func SetupConnection() (*MongoUsers, error) {
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

	col := client.Database("rzd").Collection("users")

	return &MongoUsers{
		CLI:   client,
		Users: col,
	}, nil
}

func cleanFunc(t *testing.T, clean bool, client *mongo.Collection, obj interface{}) {
	if clean {
		res, err := client.DeleteOne(context.Background(), obj)
		if err != nil {
			t.Fatal(err)
		}
		if res.DeletedCount == 0 {
			t.Fatal(errors.New(fmt.Sprintf("can't delete %v", obj)))
		}
	}
}

func TestMongoUsers_Create(t *testing.T) {
	users, err := SetupConnection()
	if err != nil {
		t.Fatal(err)
	}

	a := assert.New(t)
	testCases := []struct {
		name       string
		actualUser *entity.User
		expectedOk bool
		cleanFlag  bool
		checkErr   func(err error, msgAndArgs ...interface{}) bool
	}{
		{
			name: "first user in db",
			actualUser: &entity.User{
				UserTelegramID: "test_1",
				UserName:       "test_1_test_1",
				TrainIDS: []string{
					"1", "2", "3",
				},
			},
			expectedOk: true,
			cleanFlag:  false,
			checkErr:   a.NoError,
		}, {
			name: "existing user",
			actualUser: &entity.User{
				UserTelegramID: "test_1",
				UserName:       "test_1_test_1",
				TrainIDS: []string{
					"1", "2", "3",
				},
			},
			expectedOk: true,
			cleanFlag:  true,
			checkErr:   a.Error,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOk, err := users.Create(tc.actualUser)
			tc.checkErr(err)

			a.Equal(tc.expectedOk, actualOk)

			cleanFunc(t, tc.cleanFlag, users.Users, tc.actualUser)
		})
	}
}

func TestMongoUsers_Delete(t *testing.T) {
	users, err := SetupConnection()
	if err != nil {
		t.Fatal(err)
	}

	a := assert.New(t)
	testCases := []struct {
		name       string
		actualUser *entity.User
		expectedOk bool
		cleanFlag  bool
		checkErr   func(err error, msgAndArgs ...interface{}) bool
	}{
		{
			name: "first user in db",
			actualUser: &entity.User{
				UserTelegramID: "test_1",
				UserName:       "test_1_test_1",
				TrainIDS: []string{
					"1", "2", "3",
				},
			},
			expectedOk: true,
			cleanFlag:  false,
			checkErr:   a.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := users.Users.InsertOne(context.Background(), tc.actualUser)
			if err != nil {
				t.Fatal(err)
			}
			err = users.Delete(tc.actualUser)
			tc.checkErr(err)
		})
	}
}

func TestMongoUsers_ReadOne(t *testing.T) {
	users, err := SetupConnection()
	if err != nil {
		t.Fatal(err)
	}

	a := assert.New(t)
	testCases := []struct {
		name         string
		expectedUser *entity.User
		cleanFlag    bool
		checkErr     func(err error, msgAndArgs ...interface{}) bool
	}{
		{
			name: "success",
			expectedUser: &entity.User{
				UserTelegramID: "test_1",
				UserName:       "test_1_test_1",
				TrainIDS: []string{
					"1", "2", "3",
				},
			},
			cleanFlag: true,
			checkErr:  a.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := users.Users.InsertOne(context.Background(), tc.expectedUser)
			if err != nil {
				t.Fatal(err)
			}

			actualUser, err := users.ReadOne(tc.expectedUser)
			tc.checkErr(err)

			a.Equal(tc.expectedUser, actualUser)
			cleanFunc(t, tc.cleanFlag, users.Users, tc.expectedUser)
		})
	}
}

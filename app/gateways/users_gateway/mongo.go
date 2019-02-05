package users_gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"rzd/app/entity"
	"time"
)

/*
Init ctx for every query, coz sometimes needed more waiting time.
*/

type MongoUsers struct {
	CLI   mongo.Client
	Users mongo.Collection
}

func NewMongoUsers(cli *mongo.Client) (MongoUsers, error) {
	col := cli.Database("rzd").Collection("users")
	return MongoUsers{CLI: *cli, Users: *col}, nil
}

func (m *MongoUsers) Create(user entity.User) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	filter := user
	check := entity.User{}

	err := m.Users.FindOne(ctx, filter).Decode(&check)
	if err != nil {
		return errors.New(fmt.Sprintf("Gateway->Users_Gateway->Create: Error in mgdb.Findone - %s", err))
	}

	if check.UserTelegramID != "" {
		return errors.New(fmt.Sprintf("Gateway->Users_Gateway->Create: Error - this user exist"))
	}

	result, err := m.Users.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if result.InsertedID == nil {
		return errors.New(fmt.Sprintf(fmt.Sprintf("Gateway->Users_Gateway->Create: Got empty result - %s", result.InsertedID)))
	}

	return nil
}

func (m *MongoUsers) ReadOne() (entity.User, error) {
	result := entity.User{}
	filter := bson.M{}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := m.Users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return entity.User{}, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s", err))
	}

	return result, nil
}

func (m *MongoUsers) ReadMany() ([]entity.User, error) {
	users := []entity.User{}
	user := entity.User{}
	filter := bson.M{}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	cur, err := m.Users.Find(ctx, filter) // FIXME: add filter
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadMany: Error in mgdb.Find - %s", err))
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := cur.Decode(&user)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadMany: Error in cursour.Decode - %s", err))
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadMany: Error in cursour.Err - %s", err))
	}

	return users, nil
}

func (m *MongoUsers) ReadSection(start, end int) ([]entity.User, error) {
	users := []entity.User{}
	user := entity.User{}
	filter := bson.M{}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	cur, err := m.Users.Find(ctx, filter) // FIXME: add filter
	opts := options.Find()
	sort := opts.SetSort(nil).SetLimit(10)
	m.Users.Find(ctx, filter, sort)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadMany: Error in mgdb.Find - %s", err))
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := cur.Decode(&user)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadMany: Error in cursour.Decode - %s", err))
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadMany: Error in cursour.Err - %s", err))
	}

	return users, nil
}

func (m *MongoUsers) Update(user entity.User) error {
	result := entity.User{}
	filter := user
	filter.TrainIDS = []string{}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := m.Users.FindOneAndUpdate(ctx, filter, user).Decode(&result)
	if err != nil {
		return errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s", err))
	}

	return nil
}

func (m *MongoUsers) Delete(user entity.User) error {
	result := entity.User{}
	filter := bson.M{} // FIXME: add filter

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := m.Users.FindOneAndDelete(ctx, filter).Decode(&result)
	if err != nil {
		return errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s", err))
	}

	return nil
}

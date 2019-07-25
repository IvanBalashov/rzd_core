package users_gateway

import (
	"context"
	"errors"
	"fmt"
	"rzd/app/entity"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

/*
Init ctx for every query, coz sometimes needed more waiting time.
*/

type MongoUsers struct {
	CLI   *mongo.Client
	Users *mongo.Collection
}

func NewMongoUsers(cli *mongo.Client) (MongoUsers, error) {
	col := cli.Database("rzd").Collection("users")
	return MongoUsers{CLI: cli, Users: col}, nil
}

func (m *MongoUsers) Create(user *entity.User) (bool, error) {
	filter := user
	check := entity.User{}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_ = m.Users.FindOne(ctx, filter).Decode(&check)

	if check.UserTelegramID != "" {
		return true, errors.New(fmt.Sprintf("Gateway->Users_Gateway->Create: Error - this user exist"))
	}

	result, err := m.Users.InsertOne(ctx, user)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Gateway->Users_Gateway->Create: Error - can't insert user in db - %s", err))
	}

	if result.InsertedID == nil {
		return false, errors.New(fmt.Sprintf(fmt.Sprintf("Gateway->Users_Gateway->Create: Got empty result - %s", result.InsertedID)))
	}

	return true, nil
}

func (m *MongoUsers) ReadOne(filter *entity.User) (*entity.User, error) {
	result := entity.User{}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := m.Users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s", err))
	}

	return &result, nil
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

func (m *MongoUsers) ReadSection(start, end int64) ([]*entity.User, error) {
	users := []*entity.User{}
	user := &entity.User{}
	filter := bson.M{}
	sort := options.Find().SetLimit(end)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	cur, err := m.Users.Find(ctx, filter, sort) // FIXME: add filter
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadMany: Error in mgdb.Find - %s", err))
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := cur.Decode(user)
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

func (m *MongoUsers) Update(user *entity.User) error {
	filter := user
	filter.TrainIDS = []string{}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	result := m.Users.FindOneAndUpdate(ctx, filter, user)
	if result.Err() != nil {
		return errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s", result.Err()))
	}

	return nil
}

func (m *MongoUsers) Delete(user *entity.User) error {
	filter := user
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	result := m.Users.FindOneAndDelete(ctx, filter)
	if result.Err() != nil {
		return errors.New(fmt.Sprintf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s", result.Err()))
	}

	return nil
}

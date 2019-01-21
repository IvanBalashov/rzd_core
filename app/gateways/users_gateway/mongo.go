package users_gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
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
	result, err := m.Users.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	if result.InsertedID == nil {
		err := fmt.Sprintf("Gateway->Users_Gateway->Create: Got empty result - %s\n", result.InsertedID)
		log.Printf(err)
		return errors.New(err)
	}
	return nil
}

func (m *MongoUsers) ReadOne() (entity.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	result := entity.User{}
	filter := bson.M{}
	err := m.Users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s\n", err)
		return entity.User{}, err
	}
	return result, nil
}

func (m *MongoUsers) ReadMany(ids []int) ([]entity.User, error) {
	users := []entity.User{}
	user := entity.User{}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := m.Users.Find(ctx, nil) // FIXME: add filter
	if err != nil {
		log.Printf("Gateway->Users_Gateway->ReadMany: Error in mgdb.Find - %s\n", err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		err := cur.Decode(&user)
		if err != nil {
			log.Printf("Gateway->Users_Gateway->ReadMany: Error in cursour.Decode - %s\n", err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		log.Printf("Gateway->Users_Gateway->ReadMany: Error in cursour.Err - %s\n", err)
		return nil, err
	}
	return users, nil
}

func (m *MongoUsers) Update(user entity.User) error {
	result := entity.User{}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	filter := bson.M{} // FIXME: add filter
	err := m.Users.FindOneAndUpdate(ctx, filter, user).Decode(&result)
	if err != nil {
		log.Printf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s\n", err)
		return err
	}
	return nil
}

func (m *MongoUsers) Delete(user entity.User) error {
	result := entity.User{}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	filter := bson.M{} // FIXME: add filter
	err := m.Users.FindOneAndDelete(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("Gateway->Users_Gateway->ReadOne: Error in mgdb.FindOne - %s\n", err)
		return err
	}
	return nil
}

package trains_gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"

	"rzd/app/entity"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

/*
Init ctx for every query, coz sometimes needed more waiting time.
*/
type MongoTrains struct {
	CLI    *mongo.Client
	Trains *mongo.Collection
}

func NewMongoTrains(cli *mongo.Client) (MongoTrains, error) {
	col := cli.Database("rzd").Collection("trains")
	return MongoTrains{CLI: cli, Trains: col}, nil
}

func (m *MongoTrains) Create(train *entity.Train) (string, error) {
	filter := train
	check := entity.Train{}
	ctx, _ := context.WithTimeout(context.Background(), 1000*time.Second)
	_ = m.Trains.FindOne(ctx, filter).Decode(&check)

	if check.ID != "" {
		return "", errors.New(fmt.Sprintf("Gateway->Trains_Gateway->Create: Error - this train exist"))
	}

	result, err := m.Trains.InsertOne(ctx, train)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Gateway->Trains_Gateway->Create: Error - can't insert train in db - %s", err))
	}

	if result.InsertedID == nil {
		return "", errors.New(fmt.Sprintf(fmt.Sprintf("Gateway->Trains_Gateway->Create: Got empty result - %s", result.InsertedID)))
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	} else {
		return "", errors.New("can't get oid")
	}
}

func (m *MongoTrains) ReadOne(filter *entity.Train) (*entity.Train, error) {
	result := entity.Train{}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := m.Trains.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s", err))
	}

	return &result, nil
}

func (m *MongoTrains) ReadMany() ([]*entity.Train, error) {
	trains := []*entity.Train{}
	train := &entity.Train{}
	filter := bson.M{}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	cur, err := m.Trains.Find(ctx, filter) // FIXME: add filter
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadMany: Error in mgdb.Find - %s", err))
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := cur.Decode(train)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadMany: Error in cursour.Decode - %s", err))
		}
		trains = append(trains, train)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadMany: Error in cursour.Err - %s", err))
	}

	return trains, nil
}

func (m *MongoTrains) ReadSection(start, end int64) ([]*entity.Train, error) {
	trains := []*entity.Train{}
	train := &entity.Train{}
	filter := bson.M{}
	sort := options.Find().SetLimit(end)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	cur, err := m.Trains.Find(ctx, filter, sort) // FIXME: add filter
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadMany: Error in mgdb.Find - %s", err))
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := cur.Decode(train)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadMany: Error in cursour.Decode - %s", err))
		}
		trains = append(trains, train)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadMany: Error in cursour.Err - %s", err))
	}

	return trains, nil
}

func (m *MongoTrains) Update(train *entity.Train) error {
	filter := train
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	result := m.Trains.FindOneAndUpdate(ctx, filter, train)
	if result.Err() != nil {
		return errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s", result.Err()))
	}

	return nil
}

func (m *MongoTrains) Delete(train *entity.Train) error {
	filter := train
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	result := m.Trains.FindOneAndDelete(ctx, filter)
	if result.Err() != nil {
		return errors.New(fmt.Sprintf("Gateway->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s", result.Err()))
	}

	return nil
}

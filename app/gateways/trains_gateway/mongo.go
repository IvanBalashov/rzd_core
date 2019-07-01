package trains_gateway

import (
	"context"
	"errors"
	"fmt"
	"time"

	"rzd/app/entity"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

/*
Init ctx for every query, coz sometimes needed more waiting time.
*/
type MongoTrains struct {
	CLI    mongo.Client
	Trains mongo.Collection
}

func NewMongoTrains(cli *mongo.Client) (MongoTrains, error) {
	col := cli.Database("rzd").Collection("trains")
	return MongoTrains{CLI: *cli, Trains: *col}, nil
}

func (m *MongoTrains) Create(train entity.Train) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	result, err := m.Trains.InsertOne(ctx, train)
	if err != nil {
		return "", errors.New(fmt.Sprintf("MDB:Gateways->Trains_Gateway->Create: Error in mgdb.InsertOne - %s", err))
	}

	if result.InsertedID == nil {
		return "", errors.New(fmt.Sprintf("MDB:Gateways->Trains_Gateway->Create: Got empty result - %s", result.InsertedID))
	}

	trainID := result.InsertedID.(string)

	return trainID, nil
}

func (m *MongoTrains) ReadOne(trainID string) (entity.Train, error) {
	result := entity.Train{}
	filter := bson.M{"_id": trainID}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := m.Trains.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return entity.Train{},
			errors.New(fmt.Sprintf("MDB:Gateways->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s", err))
	}

	return result, nil
}

func (m *MongoTrains) ReadMany() ([]entity.Train, error) {
	trains := []entity.Train{}
	train := entity.Train{}
	filter := bson.M{}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	cur, err := m.Trains.Find(ctx, filter) // FIXME: add filter
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("MDB:Gateways->Trains_Gateway->ReadMany: Error in mgdb.Find - %s %s", err, cur))
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		err := cur.Decode(&train)
		if err != nil {
			return nil,
				errors.New(fmt.Sprintf("MDB:Gateways->Trains_Gateway->ReadMany: Error in cursour.Decode - %s", err))
		}
		trains = append(trains, train)
	}

	if err := cur.Err(); err != nil {
		return nil,
			errors.New(fmt.Sprintf("MDB:Gateways->Trains_Gateway->ReadMany: Error in cursour.Err - %s", err))
	}

	return trains, nil
}

func (m *MongoTrains) Update(train entity.Train) error {
	result := entity.Train{}
	filter := bson.M{} // FIXME: add filter

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := m.Trains.FindOneAndUpdate(ctx, filter, train).Decode(&result)
	if err != nil {
		return errors.New(fmt.Sprintf("MDB:Gateways->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s", err))
	}

	return nil
}

func (m *MongoTrains) Delete(train entity.Train) error {
	result := entity.Train{}
	filter := bson.M{} // FIXME: add filter

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := m.Trains.FindOneAndDelete(ctx, filter).Decode(&result)
	if err != nil {
		return errors.New(fmt.Sprintf("MDB:Gateways->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s", err))
	}

	return nil
}

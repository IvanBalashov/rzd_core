package trains_gateway

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
type MongoTrains struct {
	CLI    mongo.Client
	Trains mongo.Collection
}

func NewMongoTrains(cli *mongo.Client) (MongoTrains, error) {
	col := cli.Database("rzd").Collection("trains")
	return MongoTrains{CLI: *cli, Trains: *col}, nil
}

func (m *MongoTrains) Create(train entity.Train) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	result, err := m.Trains.InsertOne(ctx, train)
	if err != nil {
		return err
	}
	if result.InsertedID == nil {
		err := fmt.Sprintf("MDB:Gateways->Trains_Gateway->Create: Got empty result - %s\n", result.InsertedID)
		log.Printf(err)
		return errors.New(err)
	}
	return nil
}

func (m *MongoTrains) ReadOne() (entity.Train, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	result := entity.Train{}
	filter := bson.M{}
	err := m.Trains.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("MDB:Gateways->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s\n", err)
		return entity.Train{}, err
	}
	return result, nil
}

func (m *MongoTrains) ReadMany(ids []int) ([]entity.Train, error) {
	trains := []entity.Train{}
	train := entity.Train{}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := m.Trains.Find(ctx, nil) // FIXME: add filter
	if err != nil {
		log.Printf("MDB:Gateways->Trains_Gateway->ReadMany: Error in mgdb.Find - %s\n", err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		err := cur.Decode(&train)
		if err != nil {
			log.Printf("MDB:Gateways->Trains_Gateway->ReadMany: Error in cursour.Decode - %s\n", err)
			return nil, err
		}
		trains = append(trains, train)
	}
	if err := cur.Err(); err != nil {
		log.Printf("MDB:Gateways->Trains_Gateway->ReadMany: Error in cursour.Err - %s\n", err)
		return nil, err
	}
	return trains, nil
}

func (m *MongoTrains) Update(train entity.Train) error {
	result := entity.Train{}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	filter := bson.M{} // FIXME: add filter
	err := m.Trains.FindOneAndUpdate(ctx, filter, train).Decode(&result)
	if err != nil {
		log.Printf("MDB:Gateways->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s\n", err)
		return err
	}
	return nil
}

func (m *MongoTrains) Delete(train entity.Train) error {
	result := entity.Train{}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	filter := bson.M{} // FIXME: add filter
	err := m.Trains.FindOneAndDelete(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("MDB:Gateways->Trains_Gateway->ReadOne: Error in mgdb.FindOne - %s\n", err)
		return err
	}
	return nil
}

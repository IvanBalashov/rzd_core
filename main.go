package main

import (
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"os"
	"rzd/app/gateways/rzd_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"rzd/app/usecase"
	"rzd/server/http"
	"rzd/server/rabbitmq"
	"time"
)

type Config struct {
	HttpHost    string
	HttpPort    string
	PostgresUrl string
	RabbitMQUrl string
	MongoDBUrl  string
}

func GenConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Main->GenConfig: Error while load .env file - %s\n", err.Error())
	} else {
		log.Println("Main->GenConfig: File .env loaded")
	}
	var conf = Config{}
	if val, ok := os.LookupEnv("HTTP_HOST"); !ok {
		log.Printf("Main->GenConfig: HTTP_HOST env don't seted\n")
		os.Exit(2)
	} else {
		conf.HttpHost = val
	}
	if val, ok := os.LookupEnv("HTTP_PORT"); !ok {
		log.Printf("Main->GenConfig: HTTP_PORT env don't seted\n")
		os.Exit(2)
	} else {
		conf.HttpPort = val
	}
	if val, ok := os.LookupEnv("POSTGRES_URL"); !ok {
		log.Printf("Main->GenConfig: POSTGRES_URL env don't seted\n")
		os.Exit(2)
	} else {
		conf.PostgresUrl = val
	}
	if val, ok := os.LookupEnv("RABBITMQ_URL"); !ok {
		log.Printf("Main->GenConfig: RABBITMQ_URL env don't seted\n")
		os.Exit(2)
	} else {
		conf.RabbitMQUrl = val
	}
	if val, ok := os.LookupEnv("MONGODB_URL"); !ok {
		log.Printf("Main->GenConfig: MONGODB_URL env don't seted\n")
		os.Exit(2)
	} else {
		conf.MongoDBUrl = val
	}

	return conf
}

func main() {
	log.SetFlags(log.LstdFlags)

	config := GenConfig()

	log.Printf("Main: Starting app.\n")
	log.Printf("Main: Init rzd.ru REST api.\n")

	CLI := rzd_gateway.NewRestAPIClient(
		"https://pass.rzd.ru/timetable/public/ru",
		"http://www.rzd.ru/suggester",
		5827,
		5764,
		5804,
	)

	log.Printf("Main: Success.\n")
	log.Printf("Main: Init MongoDB client.\n")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, config.MongoDBUrl)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Main: Can't connect to Mongodb - %s\n", err)
	}

	MDDBTrains, err := trains_gateway.NewMongoTrains(client)
	if err != nil {
		log.Printf("Main: Can't connect to train collections - %s\n", err)
		return
	}
	MDDBUsers, err := users_gateway.NewMongoUsers(client)
	if err != nil {
		log.Printf("Main: Can't connect to users collections - %s\n", err)
		return
	}
	log.Printf("Main: Success.\n")

	app := usecase.NewApp(&MDDBTrains, &MDDBUsers, &CLI)

	// RabbitMQ Server
	{
		server, err := rabbitmq.NewServer(config.RabbitMQUrl, &app)
		if err != nil {
			log.Printf("Main: Can't connect to rabbitmq on addr - %s\n", config.RabbitMQUrl)
		} else {
			// TODO: Remove after complete rabbitmq files.
			// TODO: Think about call to another nodes about starting??
			log.Printf("Main: Success\n")
			request := rabbitmq.NewRequestQueue(&server.Chanel,
				"test",
				"",
				false,
				false,
				false,
				false,
				nil)

			response := rabbitmq.NewResponseQueue(&server.Chanel,
				"test",
				"",
				false,
				false,
				false,
				false,
				nil)

			go server.Serve(request, response)
			msg := rabbitmq.Message{
				Event: "Get",
				Data:  "kek",
			}
			time.Sleep(time.Second)
			data, _ := json.Marshal(msg)
			err := response.Send(data)
			if err != nil {
				log.Printf("Main: Error in test send - %s\n", err.Error())
			}
		}
	}
	// REST Server.
	{
		log.Printf("Main: Starting web server on addr - %s:%s", config.HttpHost, config.HttpPort)
		server := http.NewServer(http.NewHandler(&app), config.HttpHost, config.HttpPort)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Main: Error while serving - \n\t%s\n", err.Error())
			return
		}
	}
}

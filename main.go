package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"os"
	"rzd/app/gateways/rzd_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"rzd/app/usecase"
	"rzd/reporting"
	"rzd/server/http"
	"rzd/server/rabbitmq"
	"time"
)

type Config struct {
	AppName     string
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
	if val, ok := os.LookupEnv("APP_NAME"); !ok {
		log.Printf("Main->GenConfig: APP_NAME env don't seted\n")
		os.Exit(2)
	} else {
		conf.AppName = val
	}
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
	config := GenConfig()

	logs := make(chan string)
	logger := reporting.NewLogger(logs, config.AppName)
	logger.Start()

	logger.Write("Main: Starting app.")
	logger.Write("Main: Init rzd.ru REST api.")

	CLI := rzd_gateway.NewRestAPIClient(
		"https://pass.rzd.ru/timetable/public/ru",
		"http://www.rzd.ru/suggester",
		5827,
		5764,
		5804,
	)

	logger.Write("Main: Success.")
	logger.Write("Main: Init MongoDB client.")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, config.MongoDBUrl)

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Write(fmt.Sprintf("Main: Can't connect to Mongodb - %s", err.Error()))
	}

	MDDBTrains, err := trains_gateway.NewMongoTrains(client)
	if err != nil {
		logger.Write(fmt.Sprintf("Main: Can't connect to train collections - %s", err.Error()))
		return
	}
	MDDBUsers, err := users_gateway.NewMongoUsers(client)
	if err != nil {
		logger.Write(fmt.Sprintf("Main: Can't connect to users collections - %s", err.Error()))
		return
	}
	logger.Write("Main: Success.")

	app := usecase.NewApp(&MDDBTrains, &MDDBUsers, &CLI, logs)

	// RabbitMQ Server
	{
		server, err := rabbitmq.NewServer(config.RabbitMQUrl, &app)
		if err != nil {
			logger.Write(fmt.Sprintf("Main: Can't connect to rabbitmq on addr - %s", config.RabbitMQUrl))
		} else {
			// TODO: Remove after complete rabbitmq files.
			// TODO: Think about call to another nodes about starting??
			logger.Write("Main: Success")
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
				logger.Write(fmt.Sprintf("Main: Error in test send - %s", err.Error()))
			}
		}
	}
	// REST Server.
	{
		logger.Write(fmt.Sprintf("Main: Starting web server on addr - %s:%s", config.HttpHost, config.HttpPort))
		server := http.NewServer(http.NewHandler(&app), config.HttpHost, config.HttpPort)
		if err := server.ListenAndServe(); err != nil {
			logger.Write(fmt.Sprintf("Main: Error while serving - \n\t%s", err.Error()))
			return
		}
	}
}

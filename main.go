package main

import (
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"os"
	"rzd/app/gateways/route_gateway"
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
		log.Println("Error while load .env file - %s\n", err.Error())
	} else {
		log.Println("File .env loaded")
	}
	var conf = Config{}
	if val, ok := os.LookupEnv("HTTP_HOST"); !ok {
		log.Printf("HTTP_HOST env don't seted\n")
		os.Exit(2)
	} else {
		conf.HttpHost = val
	}
	if val, ok := os.LookupEnv("HTTP_PORT"); !ok {
		log.Printf("HTTP_PORT env don't seted\n")
		os.Exit(2)
	} else {
		conf.HttpPort = val
	}
	if val, ok := os.LookupEnv("POSTGRES_URL"); !ok {
		log.Printf("POSTGRES_URL env don't seted\n")
		os.Exit(2)
	} else {
		conf.PostgresUrl = val
	}
	if val, ok := os.LookupEnv("RABBITMQ_URL"); !ok {
		log.Printf("RABBITMQ_URL env don't seted\n")
		os.Exit(2)
	} else {
		conf.RabbitMQUrl = val
	}

	if val, ok := os.LookupEnv("MONGODB_URL"); !ok {
		log.Printf("RABBITMQ_URL env don't seted\n")
		os.Exit(2)
	} else {
		conf.MongoDBUrl = val
	}

	return conf
}

func main() {
	log.SetFlags(log.LstdFlags)

	config := GenConfig()

	log.Printf("Starting app.\n")
	log.Printf("Init rzd.ru REST api.\n")

	CLI := route_gateway.NewRestAPIClient(
		"https://pass.rzd.ru/timetable/public/ru",
		5827,
		5764,
		5804,
	)

	/*	log.Printf("Success.\n")
		log.Printf("Init postgres client.\n")
		connect, err := sqlx.Connect("postgres", config.PostgresUrl)
		if err != nil {
			log.Printf("Error while connect to PostgreSQL - %s\n", err)
			return
		}
		log.Printf("Success\n")
	*/
	log.Printf("Success.\n")
	log.Printf("Init MongoDB client.\n")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, config.MongoDBUrl)

	MDDBTrains, err := trains_gateway.NewMongoTrains(client)
	if err != nil {
		log.Printf("Main: Can't connect to Mongodb - %s\n", err)
		return
	}
	MDDBUsers, err := users_gateway.NewMongoUsers(client)
	if err != nil {
		log.Printf("Main: Can't connect to Mongodb - %s\n", err)
		return
	}
	log.Printf("Success.\n")

	//PGTrains := trains_gateway.NewPostgres(connect)
	//PGUsers := users_gateway.NewPostgres(connect)
	app := usecase.NewApp(&MDDBTrains, &MDDBUsers, &CLI)

	// RabbitMQ Server
	{
		server, err := rabbitmq.NewServer(config.RabbitMQUrl, &app)
		if err != nil {
			log.Printf("Can't connect to rabbitmq on addr - %s\n", config.RabbitMQUrl)
		} else {
			// TODO: Remove after complete rabbitmq files.
			// TODO: Think about call to another nodes about starting??
			log.Printf("Success\n")
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
				log.Printf("Error in test send - %s\n", err.Error())
			}
		}
	}
	// REST Server.
	{
		log.Printf("Starting web server on addr - %s:%s", config.HttpHost, config.HttpPort)
		server := http.NewServer(http.NewHandler(&app), config.HttpHost, config.HttpPort)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Error while serving - \n\t%s\n", err.Error())
			return
		}
	}
}

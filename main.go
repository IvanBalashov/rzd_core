package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"os"
	"rzd/app/gateways/cache_gateway"
	"rzd/app/gateways/rzd_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"rzd/app/usecase"
	"rzd/reporting"
	"rzd/server/http"
	"rzd/server/rabbitmq"
	"rzd/server/rabbitmq/middleware"
	"time"
)

type Config struct {
	AppName     string
	HttpHost    string
	HttpPort    string
	PostgresUrl string
	RabbitMQUrl string
	MongoDBUrl  string
	MemcacheUrl string
}

func init() {
	log.SetFlags(log.LstdFlags)
}

func GenConfig() Config {
	err := godotenv.Load()
	var appName string

	if err != nil {
		log.Printf("%s__Main->GenConfig: Error while load .env file - %s\n", os.Getenv("APP_NAME"), err.Error())
	} else {
		log.Printf("%s__Main->GenConfig: File .env loaded\n", os.Getenv("APP_NAME"))
	}
	var conf = Config{}
	if val, ok := os.LookupEnv("APP_NAME"); !ok {
		log.Printf("(Can't get app name)__Main->GenConfig: APP_NAME env don't seted\n")
		os.Exit(2)
	} else {
		conf.AppName = val
		appName = val
	}
	if val, ok := os.LookupEnv("HTTP_HOST"); !ok {
		log.Printf("%s__Main->GenConfig: HTTP_HOST env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.HttpHost = val
	}
	if val, ok := os.LookupEnv("HTTP_PORT"); !ok {
		log.Printf("%s__Main->GenConfig: HTTP_PORT env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.HttpPort = val
	}
	if val, ok := os.LookupEnv("POSTGRES_URL"); !ok {
		log.Printf("%s__Main->GenConfig: POSTGRES_URL env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.PostgresUrl = val
	}
	if val, ok := os.LookupEnv("RABBITMQ_URL"); !ok {
		log.Printf("%s__Main->GenConfig: RABBITMQ_URL env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.RabbitMQUrl = val
	}
	if val, ok := os.LookupEnv("MONGODB_URL"); !ok {
		log.Printf("%s__Main->GenConfig: MONGODB_URL env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.MongoDBUrl = val
	}
	if val, ok := os.LookupEnv("MEMCACHE_URL"); !ok {
		log.Printf("%s__Main->GenConfig: MONGODB_URL env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.MemcacheUrl = val
	}

	return conf
}

func main() {
	config := GenConfig()

	logs := make(chan string)
	defer close(logs)

	logger := reporting.NewLogger(logs, config.AppName)
	logger.Start()

	time.Sleep(100 * time.Millisecond)

	logs <- fmt.Sprintf("Main: Starting app")
	logs <- fmt.Sprintf("Main: Init rzd.ru REST api")

	CLI := rzd_gateway.NewRestAPIClient(
		"https://pass.rzd.ru/timetable/public/ru",
		"http://www.rzd.ru/suggester",
		5827,
		5764,
		5804,
	)
	logs <- fmt.Sprintf("Main: Success")

	logs <- fmt.Sprintf("Main: Connecting to MongoDB on addr - %s", config.MongoDBUrl)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, config.MongoDBUrl)
	if err != nil {
		logs <- fmt.Sprintf("Main: Can't create client of Mongodb - %s", err.Error())
		time.Sleep(500 * time.Millisecond)
		os.Exit(2)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logs <- fmt.Sprintf("Main: Can't connect to Mongodb - %s", err.Error())
		time.Sleep(500 * time.Millisecond)
		os.Exit(2)
	}

	MDDBTrains, err := trains_gateway.NewMongoTrains(client)
	if err != nil {
		logs <- fmt.Sprintf("Main: Can't connect to train collections - %s", err.Error())
		time.Sleep(500 * time.Millisecond)
		os.Exit(2)
	}
	MDDBUsers, err := users_gateway.NewMongoUsers(client)
	if err != nil {
		logs <- fmt.Sprintf("Main: Can't connect to users collections - %s", err.Error())
		time.Sleep(500 * time.Millisecond)
		os.Exit(2)
	}
	logs <- fmt.Sprintf("Main: Success")

	logs <- fmt.Sprintf("Main: Connecting to Memcache on addr - %s", config.MemcacheUrl)

	cacheCLI := memcache.New(config.MemcacheUrl)
	cache := cache_gateway.NewMemcache(*cacheCLI, 60)

	logs <- fmt.Sprintf("Main: Success")

	app := usecase.NewApp(&MDDBTrains, &MDDBUsers, &CLI, &cache, logs)

	// RabbitMQ Server
	{
		logs <- fmt.Sprintf("Main: Connecting to rabbitMQ on addr - %s", config.RabbitMQUrl)
		server, err := rabbitmq.NewServer(config.RabbitMQUrl, &app, logs)
		if err != nil {
			logs <- fmt.Sprintf("Main: Can't connect to rabbitmq on addr - %s", config.RabbitMQUrl)
			time.Sleep(500 * time.Millisecond)
			os.Exit(2)
		} else {
			// TODO: Remove after complete rabbitmq files.
			// TODO: Think about call to another nodes about starting??
			logs <- fmt.Sprintf("Main: Success.")
			trainsRequest := rabbitmq.NewRequestQueue(server.Chanel,
				"Get_all_trains",
				"",
				false,
				false,
				false,
				false,
				nil)

			trainsResponse := rabbitmq.NewResponseQueue(server.Chanel,
				"Send_all_trains",
				"",
				false,
				false,
				false,
				false,
				nil)
			// only for testing
			testResponse := rabbitmq.NewResponseQueue(server.Chanel,
				"Get_all_trains",
				"",
				false,
				false,
				false,
				false,
				nil)

			go server.Serve(trainsRequest, trainsResponse)
			msg := rabbitmq.MessageRabbitMQ{
				ID:    1,
				Event: "trains_list",
				Data: middleware.AllTrainsRequest{
					Direction: "0",
					Target:    "Москва",
					Source:    "Ярославль",
					Date:      "11.02.2019",
				},
			}
			time.Sleep(time.Second)
			data, _ := json.Marshal(msg)
			err := testResponse.Send(data)
			if err != nil {
				logs <- fmt.Sprintf("Main: Error in test send - %s", err.Error())
			}
		}
	}
	go app.Run("60s")
	// REST Server.
	{
		logs <- fmt.Sprintf("Main: Starting web server on addr - %s:%s", config.HttpHost, config.HttpPort)
		server := http.NewServer(http.NewHandler(&app), config.HttpHost, config.HttpPort)
		if err := server.ListenAndServe(); err != nil {
			logs <- fmt.Sprintf("Main: Error while serving - \n\t%s", err.Error())
			time.Sleep(500 * time.Millisecond)
			os.Exit(2)
		}
	}
}

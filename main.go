package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Config struct {
	AppName     string
	HTTPHost    string
	HTTPPort    string
	PostgresURL string
	RabbitMQURL string
	MongoDBURL  string
	MemcacheURL string
	MemcacheTTL int64
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
		conf.HTTPHost = val
	}
	if val, ok := os.LookupEnv("HTTP_PORT"); !ok {
		log.Printf("%s__Main->GenConfig: HTTP_PORT env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.HTTPPort = val
	}
	if val, ok := os.LookupEnv("RABBITMQ_URL"); !ok {
		log.Printf("%s__Main->GenConfig: RABBITMQ_URL env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.RabbitMQURL = val
	}
	if val, ok := os.LookupEnv("MONGODB_URL"); !ok {
		log.Printf("%s__Main->GenConfig: MONGODB_URL env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.MongoDBURL = val
	}
	if val, ok := os.LookupEnv("MEMCACHE_URL"); !ok {
		log.Printf("%s__Main->GenConfig: MONGODB_URL env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.MemcacheURL = val
	}
	if val, ok := os.LookupEnv("MEMCACHE_TTL"); !ok {
		log.Printf("%s__Main->GenConfig: MONGODB_URL env don't seted\n", appName)
		os.Exit(2)
	} else {
		conf.MemcacheTTL, _ = strconv.ParseInt(val, 10, 64)
	}
	return conf
}

func main() {
	config := GenConfig()

	logs := make(chan string)
	defer close(logs)

	logger := reporting.NewLogger(logs, config.AppName)
	logger.Start()

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

	logs <- fmt.Sprintf("Main: Connecting to MongoDB on addr - %s", config.MongoDBURL)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, config.MongoDBURL)
	if err != nil {
		logs <- fmt.Sprintf("Main: Can't create client of Mongodb - fail in connect to base \n\t\t err - %s", err.Error())
		time.Sleep(500 * time.Millisecond)
		os.Exit(2)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logs <- fmt.Sprintf("Main: Can't connect to Mongodb - fail to ping base \n\t\t err - %s", err.Error())
		time.Sleep(500 * time.Millisecond)
		os.Exit(2)
	}

	MDDBTrains, err := trains_gateway.NewMongoTrains(client)
	if err != nil {
		logs <- fmt.Sprintf("Main: Can't connect to train collections - MDDBTrains\n\t\t err -  %s", err.Error())
		time.Sleep(500 * time.Millisecond)
		os.Exit(2)
	}
	MDDBUsers, err := users_gateway.NewMongoUsers(client)
	if err != nil {
		logs <- fmt.Sprintf("Main: Can't connect to users collections - MDDBUsers\n\t\t err - %s", err.Error())
		time.Sleep(500 * time.Millisecond)
		os.Exit(2)
	}
	logs <- fmt.Sprintf("Main: Success")

	logs <- fmt.Sprintf("Main: Connecting to Memcache on addr - %s", config.MemcacheURL)

	cacheCLI := memcache.New(config.MemcacheURL)
	cache := cache_gateway.NewMemcache(cacheCLI, int32(config.MemcacheTTL))

	logs <- fmt.Sprintf("Main: Success")

	app := usecase.NewApp(&MDDBTrains, &MDDBUsers, &CLI, cache, logs)

	// RabbitMQ Server
	{
		logs <- fmt.Sprintf("Main: Connecting to rabbitMQ on addr - %s", config.RabbitMQURL)
		server, err := rabbitmq.NewServer(config.RabbitMQURL, &app, logs)
		if err != nil {
			logs <- fmt.Sprintf("Main: Can't connect to rabbitmq on addr - %s\n\t\t err - %s", config.RabbitMQURL, err.Error())
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
					Date:      "27.07.2019",
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
		logs <- fmt.Sprintf("Main: Starting web server on addr - %s:%s", config.HTTPHost, config.HTTPPort)
		server := http.NewServer(http.NewHandler(&app), config.HTTPHost, config.HTTPPort)
		if err := server.ListenAndServe(); err != nil {
			logs <- fmt.Sprintf("Main: Error while serving - \n\t%s", err.Error())
			time.Sleep(500 * time.Millisecond)
			os.Exit(2)
		}
	}
}

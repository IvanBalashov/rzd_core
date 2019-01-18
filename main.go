package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
	"rzd/app/gateways/route_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"rzd/app/usecase"
	"rzd/server/http"
)

type Config struct {
	HttpHost string
	HttpPort string
	PostgresUrl string
	RabbitMQHost string
	RabbitMQPort string
}

func GenConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(".env loaded.")
	}
	var conf = Config{}
	if val, ok := os.LookupEnv("HTTP_HOST"); !ok {
		log.Printf("HTTP_HOST env don't seted.\n")
		os.Exit(2)
	} else {
		conf.HttpHost = val
	}
	if val, ok := os.LookupEnv("HTTP_PORT"); !ok {
		log.Printf("HTTP_PORT env don't seted.\n")
		os.Exit(2)
	} else {
		conf.HttpPort = val
	}
	if val, ok := os.LookupEnv("POSTGRES_URL"); !ok {
		log.Printf("POSTGRES_URL env don't seted.\n")
		os.Exit(2)
	} else {
		conf.PostgresUrl = val
	}

	// TODO: add check envs for rabbitmq.
	return conf
}

func main() {
	config := GenConfig()

	log.Printf("Starting app.\n")
	log.Printf("Init rzd.ru api.\n")

	CLI := route_gateway.NewRestAPIClient(
		"https://pass.rzd.ru/timetable/public/ru",
		5827,
		5764,
		5804,
	)
	log.Printf("Success.\nInit postgres CLIent.\n")
	connect, err := sqlx.Connect("postgres", config.PostgresUrl)
	if err != nil {
		log.Printf("error - %s\n", err)
		return
	}

	PGTrains := trains_gateway.NewPostgres(connect)
	PGUsers := users_gateway.NewPostgres(connect)
	app := usecase.NewApp(&PGTrains, &PGUsers, &CLI)

	log.Printf("Succes.\nStarting web server on %s:%s", config.HttpHost, config.HttpPort)
	server := http.NewServer(http.NewHandler(&app), config.HttpHost, config.HttpPort)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("error while serving - \n\t%s\n", err.Error())
		return
	}
}

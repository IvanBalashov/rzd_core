package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"rzd/app/gateways/route_gateway"
	"rzd/app/gateways/trains_gateway"
	"rzd/app/gateways/users_gateway"
	"rzd/app/usecase"
	"rzd/server/http"
)

func main() {
	fmt.Printf("Starting app.\n")
	fmt.Printf("Init rzd.ru api.\n")
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")
	Cli := route_gateway.NewRestAPIClient(
		"https://pass.rzd.ru/timetable/public/ru",
		5827,
		5764,
		5804,
	)
	fmt.Printf("Success.\nInit postgres client.\n")
	connect, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Printf("error - %s\n", err)
		return
	}

	PGTrains := trains_gateway.NewPostgres(connect)
	PGUsers := users_gateway.NewPostgres(connect)
	app := usecase.NewApp(&PGTrains, &PGUsers, &Cli)

	fmt.Printf("Succes.\nStarting web server on %s:%s", host, port)
	server := http.NewServer(http.NewHandler(&app), host, port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("error while serving - \n\t%s\n", err.Error())
		return
	}
}

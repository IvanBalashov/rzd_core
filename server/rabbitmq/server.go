package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"rzd/app/usecase"
	"rzd/server/rabbitmq/middleware"
)

type MessageRabbitMQ struct {
	ID    int         `json:"id"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type RabbitServer struct {
	EventLayer middleware.EventLayer
	Connection *amqp.Connection
	Chanel     *amqp.Channel
	LogChanel  chan string
}

//Create new connection and chanel to rabbitmq.
// FIXME: Don't forgot close channel.
func NewServer(uri string, app usecase.Usecase, logChanel chan string) (RabbitServer, error) {
	connection, err := amqp.Dial(uri)
	if err != nil {
		return RabbitServer{}, err
	}

	ch, err := connection.Channel()
	if err != nil {
		return RabbitServer{}, err
	}

	return RabbitServer{
		Connection: connection,
		LogChanel:  logChanel,
		Chanel:     ch,
		EventLayer: middleware.NewEventLayer(app, logChanel),
	}, nil
}

func (r *RabbitServer) Serve(request RequestQueue, response ResponseQueue) {
	msg     := MessageRabbitMQ{}
	resp    := MessageRabbitMQ{}
	forever := make(chan bool) // FIXME: add exit statement

	requests, err := request.Read()
	if err != nil {
		r.LogChanel <- fmt.Sprintf("RabbitMQ: Start reading error - %s", err.Error())
		return
	}

	r.LogChanel <- fmt.Sprintf("RabbitMQ: Start reading messages")
	go func() {
		for request := range requests {
			r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Got message - %s", request.Body)

			err := json.Unmarshal(request.Body, &msg)
			if err != nil {
				r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Error while parse message - %s", err)
			}

			switch msg.Event {
			case "trains_list":
				answer, err := r.EventLayer.GetAllTrains(msg.Data)
				if err != nil {
					r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Error in middleware.GetInfoAboutTrains %s", err.Error())
				}

				resp = MessageRabbitMQ{
					ID:    msg.ID,
					Event: "trains_list_answer",
					Data:  answer,
				}
			case "save_one_train":
				answer, err := r.EventLayer.SaveInfoAboutTrain(msg.Data)
				if err != nil {
					r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Error in middleware.GetInfoAboutTrains %s", err.Error())
				}

				resp = MessageRabbitMQ{
					ID:    msg.ID,
					Event: "save_one_train_answer",
					Data:  answer,
				}
			case "users_count":
				answer, err := r.EventLayer.UsersCount()
				if err != nil {
					r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Error in middleware.UsersCount %s", err.Error())
				}

				resp = MessageRabbitMQ{
					ID:    msg.ID,
					Event: "users_count_answer",
					Data:  answer,
				}
			case "check_users":
				answer, err := r.EventLayer.CheckUsers(msg.Data)
				if err != nil {
					r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Error in middleware.CheckUsers %s", err.Error())
				}
				resp = MessageRabbitMQ{
					ID:    msg.ID,
					Event: "check_users_answer",
					Data:  answer,
				}
			case "exit":
				close(forever)
				break
			}
			r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Sending message - %+v", resp)
			data, err := json.Marshal(resp)
			if err != nil {
				r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Got error while parse answer - %s", err.Error())
			}

			err = response.Send(data)
			if err != nil {
				r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Got error while sending message - %s", err.Error())
			}
		}
	}()
	for {
		select {
		case _, ok := <-forever:
			if !ok {
				if err := request.Channel.Close(); err != nil {
					r.LogChanel <- fmt.Sprintf("RabbitMQ->Server: Can't close request channel - %s", err.Error())
				}
				return
			}
		}
	}
}

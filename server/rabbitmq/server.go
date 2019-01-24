package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"rzd/app/usecase"
	"rzd/server/rabbitmq/middleware"
)

type RabbitServer struct {
	Middlewares middleware.EventLayer
	Connection  amqp.Connection
	Chanel      amqp.Channel
	LogChanel   chan string
}

//Create new connection and chanel to rabbitmq.
// FIXME: Don't forgot close channel.
func NewServer(uri string, app usecase.Usecase, logChanel chan string) (RabbitServer, error) {
	var server = &RabbitServer{}
	connection, err := amqp.Dial(uri)
	if err != nil {
		return RabbitServer{}, err
	}

	server.Connection = *connection
	ch, err := connection.Channel()
	if err != nil {
		return RabbitServer{}, err
	}
	server.LogChanel = logChanel
	server.Chanel = *ch
	server.Middlewares = middleware.InitMiddleWares(app, logChanel)
	return *server, nil
}

// all works wirh rabbitmq now released like in man, need upgrade
func (r *RabbitServer) Serve(request RequestQueue, response ResponseQueue) {
	// listen msgs, call middlewares.
	getedMessage := middleware.Message{}
	messages, err := request.Read()
	if err != nil {
		r.LogChanel <- fmt.Sprintf("RabbitMQ: Error while start reading - %s", err.Error())
		return
	}
	r.LogChanel <- fmt.Sprintf("RabbitMQ: Start reading messages")
	forever := make(chan bool)
	// read messages
	go func() {
		for msg := range messages {
			r.LogChanel <- fmt.Sprintf("DEBUG:: RabbitMQ->MSG: %s", msg.Body)
			err := json.Unmarshal(msg.Body, &getedMessage)
			if err != nil {
				r.LogChanel <- fmt.Sprintf("RabbitMQ: err - %s", err)
			}
			switch getedMessage.Event {
			// Here write call middlewares.
			case "Get":
				r.LogChanel <- fmt.Sprintf("DEBUG:: Event.Get: Body:%s", getedMessage.Data)
				_, err := r.Middlewares.GetSeats(getedMessage.Data)
				if err != nil {
					r.LogChanel <- fmt.Sprintf("RabbitMQ: Error in GetSeats %s", err.Error())
				}
				// TODO: need create another queue for answers.
				//		err := response.Send([]byte{})
				//		if err != nil {
				//			log.Printf("%s\n", err.Error())
				//		}
			case "Set":
				r.LogChanel <- fmt.Sprintf("DEBUG:: event.Set: body:%s", getedMessage.Data)
				//		err := response.Send([]byte{})
				//		if err != nil {
				//			log.Printf("%s\n", err.Error())
				//		}
			case "Exit":
				break
			}
		}
	}()

	<-forever
	// need call this method after readed data
	//response.Send()
}

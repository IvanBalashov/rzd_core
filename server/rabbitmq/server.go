package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"rzd/app/usecase"
	"rzd/server/rabbitmq/middleware"
)

type RabbitServer struct {
	Middlewares middleware.EventLayer
	Connection  amqp.Connection
	Chanel      amqp.Channel
}

//Create new connection and chanel to rabbitmq.
// FIXME: Don't forgot close channel.
func NewServer(uri string, app usecase.Usecase) (RabbitServer, error) {
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
	server.Chanel = *ch
	server.Middlewares = middleware.InitMiddleWares(app)
	return *server, nil
}

// all works wirh rabbitmq now released like in man, need upgrade

func (r *RabbitServer) Serve(request RequestQueue, response ResponseQueue) {
	// listen msgs, call middlewares.
	getedMessage := Message{}
	messages, err := request.Read()
	if err != nil {
		log.Printf("RabbitMQ: Error while start reading - %s\n", err.Error())
		return
	}
	log.Printf("RabbitMQ: Start reading messages\n")
	forever := make(chan bool)
	// read messages
	go func() {
		for msg := range messages {
			log.Printf("DEBUG::RabbitMQ->MSG: %s\n", msg.Body)
			err := json.Unmarshal(msg.Body, &getedMessage)
			if err != nil {
				log.Printf("RabbitMQ: err - %s\n", err)
			}
			switch getedMessage.Event {
			// Here write call middlewares.
			case "Get":
				log.Printf("DEBUG:: Event.Get: Body:%s\n", getedMessage.Data)
				_, err := r.Middlewares.GetSeats(getedMessage.Data)
				if err != nil {
					log.Printf("RabbitMQ: Error in GetSeats %s\n", err.Error())
				}
				// TODO: need create another queue for answers.
				//		err := response.Send([]byte{})
				//		if err != nil {
				//			log.Printf("%s\n", err.Error())
				//		}
			case "Set":
				log.Printf("DEBUG::event.Set: body:%s\n", getedMessage.Data)
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

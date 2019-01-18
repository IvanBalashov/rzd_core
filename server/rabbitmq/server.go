package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"rzd/app/usecase"
)

type RabbitServer struct {
	App        usecase.Usecase
	Connection amqp.Connection
	Chanel     amqp.Channel
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
	server.App = app
	return *server, nil
}

// all works wirh rabbitmq now released like in man, need upgrade

func (r *RabbitServer) Serve(request RequestQueue, response ResponseQueue) {
	// listen msgs, call middlewares.
	getedMessage := Message{}
	messages, err := request.Read()
	if err != nil {
		return
	}
	// read messages
	forever := make(chan bool)

	go func() {
		for msg := range messages {
			fmt.Printf("live%s\n", msg.Body)
			err := json.Unmarshal(msg.Body, &getedMessage)
			if err != nil {
				fmt.Printf("err - %s\n", err)
			}
			switch getedMessage.Event {
			case "Get":
				fmt.Printf("event.Get: body:%s\n", getedMessage.Data)
				err := response.Send([]byte{})
				if err != nil {
					fmt.Printf("%s\n", err.Error())
				}
			case "Set":
				fmt.Printf("event.Set: body:%s\n", getedMessage.Data)
				err := response.Send([]byte{})
				if err != nil {
					fmt.Printf("%s\n", err.Error())
				}

			}
		}
	}()

	<-forever
	// need call this method after readed data
	//response.Send()
}

package rabbitmq

import "github.com/streadway/amqp"

type RabbitServer struct {

}

//Create new connection and chanel to rabbitmq.
// FIXME: Don't forgot close channel.
func NewChanel(uri string) (*amqp.Channel, error) {
	connection, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (r *RabbitServer) Serve() {
	// listen msgs, call middlewares.
}
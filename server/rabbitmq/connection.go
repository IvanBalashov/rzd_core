package rabbitmq

import "github.com/streadway/amqp"

func NewConnection(uri string) (*amqp.Connection, error) {
	connection, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	return connection, nil
}


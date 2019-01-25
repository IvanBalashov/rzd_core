package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type Queue struct {
	Name       string
	Exchange   string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

// Want save all information about queue, think is useful.
type RequestQueue struct {
	Queue   Queue
	MQueue  *amqp.Queue
	Channel *amqp.Channel
}

func NewRequestQueue(ch *amqp.Channel, name, exchange string, dur, del, exc, now bool, args map[string]interface{}) RequestQueue {
	q, err := ch.QueueDeclare(name, dur, del, exc, now, args)
	if err != nil {
		log.Printf("RabbitMQ->RequestQueue: Error while queue declare - %s\n", err)
		return RequestQueue{}
	}

	return RequestQueue{
		Queue: Queue{
			Name:       name,
			Exchange:   exchange,
			Durable:    dur,
			AutoDelete: del,
			Exclusive:  exc,
			NoWait:     now,
		},
		MQueue:  &q,
		Channel: ch,
	}
}

func (r *RequestQueue) Read() (<-chan amqp.Delivery, error) {
	// TODO: Rewrite args for consume!!!!
	messages, err := r.Channel.Consume(
		r.Queue.Name,     // queue
		r.Queue.Exchange, // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		log.Printf("RabbitMQ->RequestQueue: Error while consume messages - %s\n", err)
		return nil, err
	}

	return messages, nil
}

type ResponseQueue struct {
	Queue   Queue
	MQueue  *amqp.Queue
	Channel *amqp.Channel
}

func NewResponseQueue(ch *amqp.Channel, name, exchange string, dur, del, exc, now bool, args map[string]interface{}) ResponseQueue {
	declearedQueue, err := ch.QueueDeclare(name, dur, del, exc, now, args)
	if err != nil {
		log.Printf("RabbitMQ->ResponseQueue: Error while queue declare - %s\n", err)
		return ResponseQueue{}
	}

	return ResponseQueue{
		Queue: Queue{
			Name:       name,
			Exchange:   exchange,
			Durable:    dur,
			AutoDelete: del,
			Exclusive:  exc,
			NoWait:     now,
		},
		MQueue:  &declearedQueue,
		Channel: ch,
	}
}

func (r *ResponseQueue) Send(data []byte) error {
	// TODO: Rewrite args for publish!!!!
	err := r.Channel.Publish(
		r.Queue.Exchange, // target for messages
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		log.Printf("RabbitMQ->ResponseQueue: Error while publish messages - %s\n", err)
		return err
	}

	return nil
}

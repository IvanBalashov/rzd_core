package rabbitmq

import (
	"github.com/streadway/amqp"
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
	q, err := ch.QueueDeclare(name, dur, del, exc, now, args)
	if err != nil {
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
		MQueue:  &q,
		Channel: ch,
	}
}

func (r *ResponseQueue) Send(data []byte) error {
	err := r.Channel.Publish(
		r.Queue.Exchange, // target for messages
		r.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

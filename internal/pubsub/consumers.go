package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack = iota
	NackRequeue
	NackDiscard
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	isDurableQueue := true
	isAutoDelete := false
	isExclusive := false

	// if transient
	if simpleQueueType == 0 {
		isDurableQueue = false
		isAutoDelete = true
		isExclusive = true
	}

	queue, err := ch.QueueDeclare(queueName, isDurableQueue, isAutoDelete, isExclusive, false, amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	})
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queue.Name, key, exchange, false, nil)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	err = channel.Qos(10, 0, true)
	if err != nil {
		return fmt.Errorf("could not set QoS: %v", err)
	}
	ch, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(c <-chan amqp.Delivery) {
		defer channel.Close()
		for msg := range c {
			var msgBody T
			json.Unmarshal(msg.Body, &msgBody)
			ackType := handler(msgBody)

			if ackType == Ack {
				msg.Ack(false)
				fmt.Println("msg was: Acked")
			} else if ackType == NackRequeue {
				msg.Nack(false, true)
				fmt.Println("msg was: NackRequeued")
			} else if ackType == NackDiscard {
				msg.Nack(false, false)
				fmt.Println("msg was: NackDiscarded")
			}
		}
	}(ch)

	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {

	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		fmt.Println("ERROR at SubscribeGob:", err)
		return err
	}

	err = channel.Qos(10, 0, true)
	if err != nil {
		return fmt.Errorf("could not set QoS: %v", err)
	}
	ch, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println("ERROR at SubscribeGob:", err)
		return err
	}

	go func(c <-chan amqp.Delivery) {
		defer channel.Close()
		for msg := range c {
			var msgBody T
			msgBody, err := decode[T](msg.Body)
			if err != nil {
				// continue ?
				msg.Nack(false, false)
				fmt.Println("msg was: NackDiscarded")
			}

			ackType := handler(msgBody)

			if ackType == Ack {
				msg.Ack(false)
				fmt.Println("msg was: Acked")
			} else if ackType == NackRequeue {
				msg.Nack(false, true)
				fmt.Println("msg was: NackRequeued")
			} else if ackType == NackDiscard {
				msg.Nack(false, false)
				fmt.Println("msg was: NackDiscarded")
			}
		}
	}(ch)

	return nil
}

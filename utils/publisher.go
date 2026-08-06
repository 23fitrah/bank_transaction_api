package utils

import (
	"errors"
	"log"
	"time"
	"transaction_api/config"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(cfg *config.RabbitMQ, exchange, routing, queue string, payload interface{}) error {
	cfg.Mutex.RLock()
	ch := cfg.Channel
	confirms := cfg.Confirms
	cfg.Mutex.RUnlock()

	if ch == nil {
		return errors.New("rabbitmq channel nil")
	}

	body, err := sonic.Marshal(payload)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable = true
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	log.Printf("QUEUE : %s", queue)
	if err != nil {
		log.Printf("Gagal declare queue: %v", err.Error())

	}

	err = ch.QueueBind(
		q.Name,
		routing, // routing key, harus sama dengan yang dipakai saat publish
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Gagal bind queue: %v", err)
		return err
	}

	err = ch.Publish(
		exchange,
		routing,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
	if err != nil {
		log.Printf("Error PUBLISH : %w", err)
		return err
	}

	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			return errors.New("rabbitmq nack")
		}
		return nil

	case <-time.After(5 * time.Second):
		return errors.New("publish timeout")
	}
}

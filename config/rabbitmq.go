package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Url        string
	Exchange   string
	Conn       *amqp.Connection
	Channel    *amqp.Channel
	NotifyConn chan *amqp.Error
	NotifyChan chan *amqp.Error
	Confirms   chan amqp.Confirmation
	Mutex      sync.RWMutex
}

func NewRabbitMQ(config RabbitMQConfig) (*RabbitMQ, error) {
	r := &RabbitMQ{
		Url:      config.UrlRMQ,
		Exchange: config.Exchange,
	}

	if err := r.connectWithRetry(); err != nil {
		return nil, err
	}

	go r.handleReconnect()

	return r, nil
}

func logJSON(level string, message string, err error, extra map[string]interface{}) {
	logEntry := map[string]interface{}{
		"level":     level,
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   message,
	}

	if err != nil {
		logEntry["error"] = err.Error()
	}

	for k, v := range extra {
		logEntry[k] = v
	}

	jsonBytes, _ := sonic.Marshal(logEntry)
	fmt.Println(string(jsonBytes))
}

func (r *RabbitMQ) connectWithRetry() error {
	var err error

	for i := 1; i <= 10; i++ {
		err = r.connect()
		if err == nil {
			return nil
		}

		logJSON("error", "[RabbitMQ] connect attempt failed", err, map[string]interface{}{
			"attempt": i,
		})

		time.Sleep(time.Duration(i*2) * time.Second)
	}

	return err
}

func (r *RabbitMQ) connect() error {
	conn, err := amqp.Dial(r.Url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	err = ch.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	r.Mutex.Lock()
	r.Conn = conn
	r.Channel = ch
	r.NotifyConn = conn.NotifyClose(make(chan *amqp.Error))
	r.NotifyChan = ch.NotifyClose(make(chan *amqp.Error))
	r.Confirms = ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	r.Mutex.Unlock()

	logJSON("info", "[RabbitMQ] connected", nil, nil)

	return nil
}

func (r *RabbitMQ) handleReconnect() {
	for {
		select {
		case <-r.NotifyConn:
			r.reconnect()
		case <-r.NotifyChan:
			r.reconnect()
		}
	}
}

func (r *RabbitMQ) reconnect() {
	for {
		time.Sleep(5 * time.Second)

		if err := r.connect(); err == nil {
			logJSON("info", "[RabbitMQ] reconnected", nil, nil)
			return
		}
	}
}

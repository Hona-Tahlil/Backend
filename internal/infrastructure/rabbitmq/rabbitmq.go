package rabbitmq

import (
	"fmt"
	"hona/backend/bootstrap"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var config = bootstrap.Run().Env.RabbitMQ

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchanges    map[string]string
	queues       map[string]bool
	bindings     map[string][]string
	isConnected  bool
	closeChannel chan struct{}
	mu           sync.RWMutex
}

func NewRabbitMQ() *RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, config.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	_ = conn
	return &RabbitMQ{}
}

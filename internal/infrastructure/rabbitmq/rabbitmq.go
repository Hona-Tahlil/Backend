package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"hona/backend/bootstrap"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var config = bootstrap.Run().Env.RabbitMQ
var constants = bootstrap.Run().Constants.RabbitMQConstants
var channelNames = []string{constants.Channels.StorageUpload, constants.Channels.StorageUpload, constants.Channels.Notifications, constants.Channels.Emails}

type RabbitMQ struct {
	conn        *amqp.Connection
	channels    map[string]*amqp.Channel
	exchanges   map[string]string
	queues      map[string]bool
	bindings    map[string][]string
	isConnected bool
	stopMonitor chan struct{}
	mu          sync.RWMutex
}

func NewRabbitMQ() *RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, config.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	rmq := &RabbitMQ{
		conn:        conn,
		channels:    make(map[string]*amqp.Channel),
		exchanges:   make(map[string]string),
		queues:      make(map[string]bool),
		bindings:    make(map[string][]string),
		isConnected: true,
		stopMonitor: make(chan struct{}),
	}
	rmq.MakeChannels(conn, channelNames...)

	if err := rmq.declareExchange(constants.Exchanges.General, constants.Exchanges.TypeTopic); err != nil {
		err2 := rmq.Close()
		if err2 != nil {
			panic(err2)
		}
		log.Printf("error during declare exchange: %v", err)
		panic(err)
	}

	if err := rmq.setupDeadLetterQueue(); err != nil {
		err2 := rmq.Close()
		if err2 != nil {
			panic(err2)
		}
		log.Printf("error during declare DLQ: %v", err)
		panic(err)
	}

	rmq.MakeQueues(channelNames...)

	go rmq.monitorConnection()

	return rmq
}

func (rmq *RabbitMQ) MakeChannels(conn *amqp.Connection, channelNames ...string) {
	for _, ch := range channelNames {
		channel, err := conn.Channel()
		if err != nil {
			err2 := conn.Close()
			if err2 != nil {
				panic(err2)
			}
			panic(err)
		}
		rmq.channels[ch] = channel
	}
}

func (rmq *RabbitMQ) MakeQueues(channelNames ...string) {
	queues := []string{constants.Events.FileUpload, constants.Events.MultipleFilesUpload, constants.Events.SendNotification, constants.Events.SendEmail}
	for i, queue := range queues {
		if err := rmq.declareQueueWithDLX(queue, constants.Exchanges.DLX, channelNames[i]); err != nil {
			err2 := rmq.Close()
			if err2 != nil {
				panic(err2)
			}
			log.Printf("error during declare Queue: %v", err)
			panic(err)
		}
		if err := rmq.bindQueue(queue, constants.Exchanges.General, queue, channelNames[i]); err != nil {
			err2 := rmq.Close()
			if err2 != nil {
				panic(err2)
			}
			log.Printf("error during bind Queue: %v", err)
			panic(err)
		}
	}
}

func (rmq *RabbitMQ) declareExchange(name, exchangeType string) error {
	err := rmq.channels[constants.Channels.Emails].ExchangeDeclare(
		name,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err == nil {
		rmq.exchanges[name] = exchangeType
	}
	return err
}

func (rmq *RabbitMQ) setupDeadLetterQueue() error {
	if err := rmq.declareExchange(constants.Exchanges.DLX, constants.Exchanges.TypeFanout); err != nil {
		return err
	}

	_, err := rmq.channels[constants.Channels.DLQ].QueueDeclare(
		constants.Queues.DLQ,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	rmq.queues[constants.Queues.DLQ] = true

	return rmq.channels[constants.Channels.DLQ].QueueBind(
		constants.Queues.DLQ,
		"",
		constants.Exchanges.DLX,
		false,
		nil,
	)
}

func (rmq *RabbitMQ) Close() error {
	close(rmq.stopMonitor)

	for _, ch := range rmq.channels {
		if ch != nil {
			err := ch.Close()
			if err != nil {
				return err
			}
		}
	}

	if rmq.conn != nil {
		if err := rmq.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	return nil
}

func (rmq *RabbitMQ) monitorConnection() {
	connCloseChan := rmq.conn.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-rmq.stopMonitor:
			return
		case err := <-connCloseChan:
			if err != nil {
				rmq.mu.Lock()
				rmq.isConnected = false
				rmq.mu.Unlock()

				log.Printf("RabbitMQ connection lost: %v, attempting to reconnect...", err)

				for {
					rmq.mu.RLock()
					connected := rmq.isConnected
					rmq.mu.RUnlock()

					if connected {
						break
					}

					if err := rmq.reconnect(); err != nil {
						log.Printf("Failed to reconnect to RabbitMQ: %v, retrying in %s", err, config.RetryDelay)
						time.Sleep(config.RetryDelay)
					} else {
						log.Println("Successfully reconnected to RabbitMQ")
						connCloseChan = rmq.conn.NotifyClose(make(chan *amqp.Error))
						break
					}
				}
			}
		}
	}
}

func (rmq *RabbitMQ) declareQueueWithDLX(name, dlx, channel string) error {
	args := amqp.Table{
		constants.Headers.DeadLetter: dlx,
	}

	_, err := rmq.channels[channel].QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		args,
	)
	if err == nil {
		rmq.queues[name] = true
	}
	return err
}

func (rmq *RabbitMQ) bindQueue(queue, exchange, routingKey, channel string) error {
	err := rmq.channels[channel].QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err == nil {
		rmq.bindings[queue] = append(rmq.bindings[queue], routingKey)
	}
	return err
}

func (rmq *RabbitMQ) reconnect() error {
	if rmq.conn != nil {
		rmq.conn.Close()
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.User, config.Password, config.Host, config.Port, config.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	rmq.mu.Lock()
	defer rmq.mu.Unlock()

	rmq.MakeChannels(conn, channelNames...)

	rmq.conn = conn
	rmq.isConnected = true

	if err := rmq.declareExchange(constants.Exchanges.General, constants.Exchanges.TypeTopic); err != nil {
		err2 := rmq.Close()
		if err2 != nil {
			panic(err2)
		}
		log.Printf("error during declare exchange: %v", err)
		panic(err)
	}

	if err := rmq.setupDeadLetterQueue(); err != nil {
		err2 := rmq.Close()
		if err2 != nil {
			panic(err2)
		}
		log.Printf("error during declare DLQ: %v", err)
		panic(err)
	}

	rmq.MakeQueues(channelNames...)

	return nil
}

func (rmq *RabbitMQ) PublishMessage(queue string, message interface{}) error {
	rmq.mu.RLock()
	connected := rmq.isConnected
	channel, exists := rmq.channels[queue]
	rmq.mu.RUnlock()

	if !connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	if !exists || channel == nil {
		return fmt.Errorf("channel %s does not exist", queue)
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx := context.Background()
	err = channel.PublishWithContext(
		ctx,
		constants.Exchanges.General,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

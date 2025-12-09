package rabbitmq

import (
	"fmt"
	"hona/backend/bootstrap"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var config = bootstrap.Run().Env.RabbitMQ
var constants = bootstrap.Run().Constants.RabbitMQConstants

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
	channels := make(map[string]*amqp.Channel)
	ch1, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
		panic(err)
	}
	channels[constants.Channels.Chat] = ch1
	ch2, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
		panic(err)
	}
	channels[constants.Channels.Notifications] = ch2
	ch3, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
		panic(err)
	}
	channels[constants.Channels.Emails] = ch3
	ch4, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
		panic(err)
	}
	channels[constants.Channels.StorageUpload] = ch4

	rmq := &RabbitMQ{
		conn:        conn,
		channels:    channels,
		exchanges:   make(map[string]string),
		queues:      make(map[string]bool),
		bindings:    make(map[string][]string),
		isConnected: true,
		stopMonitor: make(chan struct{}),
	}

	if err := rmq.declareExchange(constants.Exchanges.General, constants.Exchanges.TypeTopic); err != nil {
		err := rmq.Close()
		if err != nil {
			panic(err)
		}
		log.Printf("error during declare exchange: %v", err)
		panic(err)
	}

	if err := rmq.setupDeadLetterQueue(); err != nil {
		err := rmq.Close()
		if err != nil {
			panic(err)
		}
		log.Printf("error during declare DLQ: %v", err)
		panic(err)
	}

	return &RabbitMQ{}
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

	queues := []string{constants.Events.FileUpload, constants.Events.MultipleFilesUpload, constants.Events.SendNotification, constants.Events.SendEmail}
	channels := []string{constants.Channels.StorageUpload, constants.Channels.StorageUpload, constants.Channels.Notifications, constants.Channels.Emails}

	for i, queue := range queues {
		if err := rmq.declareQueueWithDLX(queue, constants.Exchanges.DLX, channels[i]); err != nil {
			err := rmq.Close()
			if err != nil {
				panic(err)
			}
			log.Printf("error during declare Queue: %v", err)
			panic(err)
		}
		if err := rmq.bindQueue(queue, constants.Exchanges.General, queue, channels[i]); err != nil {
			err := rmq.Close()
			if err != nil {
				panic(err)
			}
			log.Printf("error during bind Queue: %v", err)
			panic(err)
		}
	}

	go rmq.monitorConnection()

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

					//if err := rmq.reconnect(); err != nil {
					//	log.Printf("Failed to reconnect to RabbitMQ: %v, retrying in %s", err, config.RetryDelay)
					//	time.Sleep(config.RetryDelay)
					//} else {
					//	log.Println("Successfully reconnected to RabbitMQ")
					//	connCloseChan = rmq.conn.NotifyClose(make(chan *amqp.Error))
					//	break
					//}
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

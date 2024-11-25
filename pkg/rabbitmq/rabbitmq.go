package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type HandlerFunc func(ctx context.Context, routingKey string, msg []byte) error

var (
	rabbitConn *amqp.Connection                 // Global RabbitMQ connection
	channels   = make(map[string]*amqp.Channel) // Map to store channels
	mutex      sync.Mutex                       // Mutex to ensure thread-safe map operations
)

type SvcRabbitMQ interface {
	InitializeConnection() *rabbitmq
	Close() error
	CreateChannel(routingKey string, handler HandlerFunc) *amqp.Channel
	SendMessage(routingKey string, message string) (err error)
}

type rabbitmq struct {
	appName        string
	amqpURL        string
	timeout        time.Duration
	connectionName string
	heartbeat      time.Duration
	exchange       string
}

func New(cfg Config) SvcRabbitMQ {
	return &rabbitmq{
		appName:        cfg.AppName,
		amqpURL:        cfg.AmqpURL,
		timeout:        cfg.Timeout,
		connectionName: cfg.ConnectionName,
		heartbeat:      cfg.Heartbeat,
		exchange:       cfg.Exchange,
	}
}

func (svc *rabbitmq) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (svc *rabbitmq) InitializeConnection() *rabbitmq {
	var err error

	// Create a custom dialer with timeout
	netDialer := &net.Dialer{
		Timeout: svc.timeout,
	}

	// Retry connection attempts with timeout and heartbeat interval
	for {
		// Set the connection properties with custom connection name and heartbeat interval
		rabbitConn, err = amqp.DialConfig(svc.amqpURL, amqp.Config{
			Dial: netDialer.Dial, // Set the dialer with the timeout
			Properties: amqp.Table{
				"connection_name": svc.connectionName, // Custom connection name
			},
			Heartbeat: svc.heartbeat, // Set the heartbeat interval
		})
		if err == nil {
			log.Printf("Successfully connected to rabbitmq")
			go svc.monitorHeartbeat(rabbitConn, svc.heartbeat)
			return svc
		}

		log.Printf("Failed to connect to RabbitMQ with connection name %s: %s. Retrying in 5 seconds...", svc.connectionName, err)
		time.Sleep(5 * time.Second) // Retry every 5 seconds
	}
}

func (svc *rabbitmq) Close() error {

	if rabbitConn == nil {
		return nil
	}

	if err := rabbitConn.Close(); err != nil {
		return err
	}

	rabbitConn = nil

	return nil
}

func (svc *rabbitmq) monitorHeartbeat(conn *amqp.Connection, heartbeatInterval time.Duration) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	// Channel to listen for close notifications
	closeChan := conn.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-ticker.C: // Every time the ticker ticks (heartbeat interval)
			// log.Println("Heartbeat received, connection is alive.")
		case err := <-closeChan:
			if err != nil {
				log.Printf("RabbitMQ connection closed: %v", err)
				panic("RabbitMQ down please conntact infra")
			}
		}
	}
}

func (svc *rabbitmq) CreateChannel(routingKey string, handler HandlerFunc) *amqp.Channel {
	mutex.Lock()
	defer mutex.Unlock()

	queueName := fmt.Sprintf("%s.%s.%s", svc.exchange, routingKey, svc.appName)

	// Check if the channel already exists
	if ch, exists := channels[queueName]; exists {
		return ch
	}

	// Check if the RabbitMQ connection is nil
	if rabbitConn == nil {
		log.Fatal("RabbitMQ connection is not initialized yet")
	}

	// Create a new channel if it doesn't exist
	ch, err := rabbitConn.Channel()
	svc.failOnError(err, "Failed to create a new channel")

	// Declare the topic exchange
	err = ch.ExchangeDeclare(
		svc.exchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	svc.failOnError(err, "Failed to declare an exchange")

	// Declare a queue for this subscriber
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	svc.failOnError(err, "Failed to declare a queue")

	// Bind the queue to the exchange with the provided routing key
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		svc.exchange, // exchange
		false,
		nil,
	)
	svc.failOnError(err, "Failed to bind a queue")

	// Consume messages from the queue
	go func() {
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		svc.failOnError(err, "Failed to register a consumer")

		for d := range msgs {
			err := handler(context.TODO(), routingKey, d.Body)
			if err != nil {
				log.Printf("Handler error: %v", err)
			}
		}
	}()

	// log.Printf(" [*] Subscriber for queue '%s' with routing key '%s' is running", queueName, routingKey)

	// Store the channel in the map
	channels[queueName] = ch
	return ch
}

func (svc *rabbitmq) SendMessage(routingKey string, message string) (err error) {
	// Create a channel
	ch, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	// Declare the exchange (for topic exchange)
	err = ch.ExchangeDeclare(
		svc.exchange, // exchange name
		"topic",      // exchange type (topic)
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
		return
	}

	// Publish the message to the exchange with the routing key
	err = ch.Publish(
		svc.exchange, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",    // message content type
			Body:        []byte(message), // message body
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return
	}

	// log.Printf("Sent message: %s", message)

	return
}

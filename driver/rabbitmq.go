package driver

import (
	"time"

	"github.com/skyapps-id/go-rabbitmq-exmaple/pkg/rabbitmq"
)

func NewRabbitMQ() {
	// Change using .env for proper deployment
	cfg := rabbitmq.Config{
		AppName:        "Apps",
		AmqpURL:        "amqp://guest:guest@localhost:5672/",
		Timeout:        5 * time.Second,
		Heartbeat:      10 * time.Second,
		ConnectionName: "consum",
		Exchange:       "internal",
	}

	rabbitmq.Init(cfg)
	// defer rabbitmq.Close()

}

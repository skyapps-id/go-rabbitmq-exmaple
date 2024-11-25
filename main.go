package main

import (
	"github.com/skyapps-id/go-rabbitmq-exmaple/driver"
	"github.com/skyapps-id/go-rabbitmq-exmaple/server/http"
	"github.com/skyapps-id/go-rabbitmq-exmaple/server/rabbitmq"
)

func main() {
	driver.NewRabbitMQ()

	rabbitmq.StartRabbitMQ()
	http.StartHTTP()

	rabbitmq.StopRabbitMQ()
}

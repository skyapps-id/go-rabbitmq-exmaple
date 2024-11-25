package rabbitmq

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/skyapps-id/go-rabbitmq-exmaple/pkg/rabbitmq"
)

func StartRabbitMQ() {
	// Put here for dependency injection

	SetupSubscriber()
}

func WaitServiceRabbitMQ() {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down RabbitMQ")
}

func StopRabbitMQ() {
	rabbitmq.Close()
}

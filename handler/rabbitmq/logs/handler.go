package logs

import (
	"context"
	"log"
)

type LogsHandler interface {
	Info(ctx context.Context, topic string, msg []byte) (err error)
	Error(ctx context.Context, topic string, msg []byte) (err error)
	Debug(ctx context.Context, topic string, msg []byte) (err error)
}

type handler struct {
}

func NewHandler() LogsHandler {
	return &handler{}
}

func (h handler) Info(ctx context.Context, topic string, msg []byte) (err error) {
	log.Printf("[Info Handler] Received message from topic '%s' with data: %s", topic, msg)
	return
}

func (h handler) Error(ctx context.Context, topic string, msg []byte) (err error) {
	log.Printf("[Error Handler] Received message from topic '%s' with data: %s", topic, msg)
	return
}

func (h handler) Debug(ctx context.Context, topic string, msg []byte) (err error) {
	log.Printf("[Debug Handler] Received message from topic '%s' with data: %s", topic, msg)
	return
}

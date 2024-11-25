package example

import (
	"fmt"
	"net/http"

	"github.com/skyapps-id/go-rabbitmq-exmaple/pkg/rabbitmq"
)

type ExampleHandler interface {
	HomeHandler(w http.ResponseWriter, r *http.Request)
	PublishHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
}

func NewHandler() ExampleHandler {
	return &handler{}
}

func (h handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Go HTTP Server!")
}

func (h handler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	rabbitmq.Publish("logs.info", "test logs...")
	fmt.Fprintln(w, "Success Publish!")
}

package http

import (
	"net/http"

	"github.com/skyapps-id/go-rabbitmq-exmaple/handler/http/example"
)

func Router() {
	// Handler
	exampleHandler := example.NewHandler()

	http.HandleFunc("/", exampleHandler.HomeHandler)
	http.HandleFunc("/publish", exampleHandler.PublishHandler)
}

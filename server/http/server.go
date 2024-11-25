package http

import (
	"log"
	"net/http"
)

func StartHTTP() {
	// Put here for dependency injection

	Router()

	port := "8080"
	log.Printf("Server is running on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("Error starting server: %s\n", err)
	}
}

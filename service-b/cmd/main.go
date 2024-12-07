package main

import (
	"lab-cloud-run/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handlers.CEPHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

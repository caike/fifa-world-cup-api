package main

import (
	"fifa-heroku/data"
	"fifa-heroku/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	data.PrintUsage()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/winners", handlers.WinnersHandler)
	http.ListenAndServe(":"+port, nil)
}

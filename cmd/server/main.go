package main

import (
	"fmt"
	"net/http"

	"github.com/Arush71/url-shortener/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/shorten", handlers.HandleShortening)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("server started listening....")
	server.ListenAndServe()
}

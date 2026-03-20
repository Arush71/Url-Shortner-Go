package main

import (
	"fmt"
	"net/http"

	"github.com/Arush71/url-shortener/internal/handlers"
	"github.com/Arush71/url-shortener/internal/shortner"
)

func main() {
	mux := http.NewServeMux()
	str := shortner.CreateStorage()
	handler := &handlers.Handler{
		Storage: str,
	}
	mux.HandleFunc("POST /api/shorten", handler.HandleShortening)
	mux.HandleFunc("GET /{code}", handler.Redirect)
	mux.HandleFunc("GET /stats/{code}", handler.Stats)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("server started listening....")
	server.ListenAndServe()
}

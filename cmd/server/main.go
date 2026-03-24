package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arush71/url-shortener/internal/cache"
	"github.com/Arush71/url-shortener/internal/db"
	"github.com/Arush71/url-shortener/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	APP_URL := os.Getenv("APP_URL")
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQuery := db.New(database)
	mux := http.NewServeMux()
	C := cache.SetupCache(dbQuery)
	go func() {
		ticker := time.NewTicker(4 * time.Second)
		for {
			<-ticker.C
			C.Flush()
		}
	}()
	handler := &handlers.Handler{
		C:      C,
		Q:      dbQuery,
		DB:     database,
		AppUrl: APP_URL,
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

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
	"github.com/Arush71/url-shortener/internal/middleware"
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
	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}
	dbQuery := db.New(database)
	mux := http.NewServeMux()
	C := cache.SetupCache(dbQuery)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
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
	ipManager := middleware.SetupIpManager()
	mux.HandleFunc("POST /api/shorten", ipManager.RateLimitMiddleware(20)(handler.HandleShortening))
	mux.HandleFunc("GET /{code}", ipManager.RateLimitMiddleware(150)(handler.Redirect))
	mux.HandleFunc("GET /api/stats/{code}", ipManager.RateLimitMiddleware(35)(handler.Stats))
	// Frontent Routes
	mux.HandleFunc("GET /stats.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/stats.html")
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})
	go ipManager.CleanUpIp()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("server started listening....")
	log.Fatal(server.ListenAndServe())
}
